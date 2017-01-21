#include "table_opob.h"

#include "libsaki/string_enum.h"
#include "libsaki/util.h"


TableOp::TableOp(TableOpOb &opOb, Who self) 
	: TableOperator(self)
	, mOpOb(opOb)
{
}

void TableOp::onActivated(Table &table)
{
	mOpOb.onActivated(mSelf, table);
}



Action makeAction(const string &actStr, const string &actArg)
{
	using AC = ActCode;

	AC act = actCodeOf(actStr.c_str());
	switch (act) {
		case AC::SWAP_OUT:
		case AC::ANKAN:
			return Action(act, T37(actArg.c_str()));
		case AC::CHII_AS_LEFT:
		case AC::CHII_AS_MIDDLE:
		case AC::CHII_AS_RIGHT:
		case AC::PON:
		case AC::KAKAN:
		case AC::IRS_CHECK:
		case AC::IRS_RIVAL:
			return Action(act, std::stoi(actArg));
		default:
			return Action(act);
	}
}

std::vector<bool> createSwapMask(const TileCount &closed,
                                 const std::vector<T37> &choices)
{
	// assume 'choices' is 34-sorted
	std::vector<bool> res;

	auto it = choices.begin();
	for (const T37 &t : tiles37::ORDER37) {
		if (it == choices.end())
			break;
		int ct = closed.ct(t);
		if (ct > 0) {
			bool val = t.looksSame(*it);
			while (ct --> 0)
				res.push_back(val);
			it += val; // consume choice if matched
		}
	}

	return res;
}

std::vector<string> createTileStrs(const std::vector<T34> &ts)
{
	std::vector<string> res;
	for (T34 t : ts)
		res.emplace_back(t.str());
	return res;
}

json createTile(const T37 &t, bool lay = false)
{
	json res;
	res["modelTileStr"] = t.str();
	res["modelLay"] = lay;
	return res;
}

json createTiles(const std::vector<T37> &ts)
{
	json res = json::array();
	for (const T37 &t : ts)
		res.emplace_back(createTile(t, false));
	return res;
}

json createBark(const M37 &m)
{
	json res;
	using T = M37::Type;
	T type = m.type();
	res["type"] = (type == T::CHII ? 1 : (type == T::PON ? 3 : 4));
	int open = m.layIndex();
	if (type != T::ANKAN)
		res["open"] = open;

	res["0"] = createTile(m[0], open == 0);
	res["1"] = createTile(m[1], open == 1);
	res["2"] = createTile(m[2], open == 2);

	if (m.isKan()) {
		res["3"] = createTile(m[3], type == T::KAKAN);
		res["isDaiminkan"] = (type == T::DAIMINKAN);
		res["isAnkan"] = (type == T::ANKAN);
		res["isKakan"] = (type == T::KAKAN);
	}

	return res;
}

json createBarks(const std::vector<M37> &ms)
{
	json list = json::array();
	for (const M37 &m : ms)
		list.emplace_back(createBark(m));
	return list;
}

template<typename T>
void rotate(T &arr)
{
	auto temp = arr[0];
	arr[0] = arr[1];
	arr[1] = arr[2];
	arr[2] = arr[3];
	arr[3] = temp;
}
                  



TableOpOb::TableOpOb()
	: mOps {
		TableOp(*this, Who(0)),
		TableOp(*this, Who(1)),
		TableOp(*this, Who(2)),
		TableOp(*this, Who(3))
	}
{
}

void TableOpOb::onActivated(Who who, Table &table)
{
    using AC = ActCode;
    const TicketFolder &tifo = table.getTicketFolder(who);

	int focusWho;
	if (tifo.can(AC::CHII_AS_LEFT)
			|| tifo.can(AC::CHII_AS_MIDDLE)
			|| tifo.can(AC::CHII_AS_RIGHT)
			|| tifo.can(AC::PON)
			|| tifo.can(AC::DAIMINKAN)
			|| tifo.can(AC::RON)) {
		focusWho = table.getFocus().who().turnFrom(who);
	} else {
		focusWho = -1;
	}

    json map;

	if (tifo.can(AC::SWAP_OUT)) {
		json mask;
		const TileCount &closed = table.getHand(who).closed();
		const auto &choices = tifo.swappables();
		map[stringOf(AC::SWAP_OUT)] = createSwapMask(closed, choices);
	}

	if (tifo.can(AC::ANKAN))
		map[stringOf(AC::ANKAN)] = createTileStrs(tifo.ankanables());

	if (tifo.can(AC::KAKAN))
		map[stringOf(AC::KAKAN)] = tifo.kakanables();

	if (tifo.can(AC::IRS_CHECK)) {
		const Girl &girl = table.getGirl(who);
		int prediceCount = girl.irsCheckCount();
		json list = json::array();
		for (int i = 0; i < prediceCount; i++) {
			const IrsCheckRow &row = girl.irsCheckRow(i);
			json rmap;
			rmap["modelMono"] = row.mono;
			rmap["modelIndent"] = row.indent;
			rmap["modelText"] = row.name;
			rmap["modelAble"] = row.able;
			rmap["modelOn"] = row.on;
			list.emplace_back(rmap);
		}
		map[stringOf(AC::IRS_CHECK)] = list;
	}

	if (tifo.can(AC::IRS_RIVAL)) {
		const Girl &girl = table.getGirl(who);
		std::vector<int> tars;
		for (int i = 0; i < 4; i++)
			if (girl.irsRivalMask()[i])
				tars.push_back(Who(i).turnFrom(who));
		map[stringOf(AC::IRS_RIVAL)] = tars;
	}

	static const AC just[] = {
		AC::PASS, AC::SPIN_OUT,
		AC::CHII_AS_LEFT, AC::CHII_AS_MIDDLE, AC::CHII_AS_RIGHT,
		AC::PON, AC::DAIMINKAN, AC::RIICHI,
		AC::RON, AC::TSUMO, AC::RYUUKYOKU,
		AC::END_TABLE, AC::NEXT_ROUND, AC::DICE, AC::IRS_CLICK
	};

    for (AC code : just)
        if (tifo.can(code))
            map[stringOf(code)] = true;

    json msg;
    msg["Type"] = "t-activated";
    msg["Action"] = map;
    msg["LastDiscarder"] = focusWho;
	peer(who.index(), msg);
}

void TableOpOb::onTableStarted(const Table &table, uint32_t seed)
{
	(void) seed;
	onPointsChanged(table);
}

void TableOpOb::onFirstDealerChoosen(Who initDealer)
{
	json msg;
	msg["Type"] = "t-first-dealer-choosen";
	for (int w = 0; w < 4; w++) {
		msg["InitDealer"] = initDealer.turnFrom(Who(w));
		peer(w, msg);
	}
}

void TableOpOb::onRoundStarted(int r, int e, Who d, 
                               bool al, int dp, uint32_t s)
{
	util::p("onRoundStarted", r, e, d.index(), 
            "al", al, "depo", dp, "seed", s);
	json msg;
	msg["Type"] = "t-round-started";
	msg["Round"] = r;
	msg["ExtraRound"] = e;
	msg["AllLast"] = al;
	msg["Deposit"] = dp;
	for (int w = 0; w < 4; w++) {
		msg["Dealer"] = d.turnFrom(Who(w));
		peer(w, msg);
	}
}

void TableOpOb::onCleaned()
{
	json msg;
	msg["Type"] = "t-cleaned";
	broad(msg);
}

void TableOpOb::onDiced(const Table &table, int die1, int die2)
{
	json msg;
	msg["Type"] = "t-diced";
	msg["Die1"] = die1;
	msg["Die2"] = die2;
	broad(msg);
}

void TableOpOb::onDealt(const Table &table)
{
	json msg;
	msg["Type"] = "t-dealt";
	for (int w = 0; w < 4; w++) {
		const auto &init = table.getHand(Who(w)).closed().t37s(true);
		msg["Init"] = createTiles(init);
		peer(w, msg);
	}
}

void TableOpOb::onFlipped(const Table &table)
{
	json msg;
	msg["Type"] = "t-flipped";
	msg["NewIndic"] = createTile(table.getMount().getDrids().back());
	broad(msg);
}

void TableOpOb::onDrawn(const Table &table, Who who)
{
	const T37 &in = table.getHand(who).drawn();
	for (int w = 0; w < 4; w++) {
		json msg;
		msg["Type"] = "t-drawn";
		msg["Who"] = who.turnFrom(Who(w));
		msg["Rinshan"] = table.duringKan();
		if (w == who.index())
			msg["Tile"] = createTile(in);
		peer(w, msg);
	}
}

void TableOpOb::onDiscarded(const Table &table, bool spin)
{
	Who discarder = table.getFocus().who();
	const T37 &out = table.getFocusTile();
	bool lay = table.lastDiscardLay();

	json msg;
	msg["Type"] = "t-discarded";
	msg["Tile"] = createTile(out, lay);
	msg["Spin"] = spin;
	for (int w = 0; w < 4; w++) {
		msg["Who"] = discarder.turnFrom(Who(w));
		peer(w, msg);
	}
}

void TableOpOb::onRiichiCalled(Who who)
{
	json msg;
	msg["Type"] = "t-riichi-called";
	for (int w = 0; w < 4; w++) {
		msg["Who"] = who.turnFrom(Who(w));
		peer(w, msg);
	}
}

void TableOpOb::onRiichiEstablished(Who who)
{
	json msg;
	msg["Type"] = "t-riichi-established";
	for (int w = 0; w < 4; w++) {
		msg["Who"] = who.turnFrom(Who(w));
		peer(w, msg);
	}
}

void TableOpOb::onBarked(const Table &table, Who who, 
                         const M37 &bark, bool spin)
{
	Who from = bark.isCpdmk() ? table.getFocus().who() : Who();

	json msg;
	msg["Type"] = "t-barked";
	msg["ActStr"] = stringOf(bark.type());
	msg["Bark"] = createBark(bark);
	msg["Spin"] = spin;
	for (int w = 0; w < 4; w++) {
		msg["Who"] = who.turnFrom(Who(w));
		msg["FromWhem"] = from.somebody() ? from.turnFrom(Who(w)) : -1;
		peer(w, msg);
	}
}

void TableOpOb::onRoundEnded(const Table &table, RoundResult result,
		                     const std::vector<Who> &openers, Who gunner,
		                     const std::vector<Form> &forms)
{
	using RR = RoundResult;

	// form and hand lists have same order as openers
	// but they don't need to be rotated since openers
	// are not rotated but changed by value
	json formsList = json::array();
	json handsList = json::array();

	for (Who who : openers) {
		const Hand &hand = table.getHand(who);

		json handMap;
		handMap["closed"] = createTiles(hand.closed().t37s(true));
		handMap["barks"] = createBarks(hand.barks());

		if (result == RR::TSUMO)
			handMap["pick"] = createTile(hand.drawn(), true);
		else if (result == RR::RON || result == RR::SCHR)
			handMap["pick"] = createTile(table.getFocusTile(), true);

		handsList.emplace_back(handMap);
	}

	for (size_t i = 0; i < forms.size(); i++) {
		const Form &form = forms[i];
		json formMap;
		formMap["spell"] = form.spell();
		formMap["charge"] = form.charge();
		formsList.emplace_back(formMap);
	}

	json msg;
	msg["Type"] = "t-round-ended";
	msg["Result"] = stringOf(result);
	msg["Hands"] = handsList;
	msg["Forms"] = formsList;
	msg["Urids"] = createTiles(table.getMount().getUrids());
	for (int w = 0; w < 4; w++) {
		msg["Openers"] = json();
		for (Who who : openers)
			msg["Openers"].push_back(who.turnFrom(Who(w)));
		msg["Gunner"] = gunner.somebody() ? gunner.turnFrom(Who(w)) : -1;
		peer(w, msg);
	}
}

void TableOpOb::onPointsChanged(const Table &table)
{
	json msg;
	msg["Type"] = "t-points-changed";
	msg["Points"] = table.getPoints();
	for (int w = 0; w < 4; w++) {
		peer(w, msg);
		rotate(msg["Points"]);
	}
}

void TableOpOb::onTableEnded(const std::array<Who, 4> &rank,
		                     const std::array<int, 4> &scores)
{
	mEnd = true;

	json msg;
	msg["Type"] = "t-table-ended";
	msg["Scores"] = scores;
	for (int w = 0; w < 4; w++) {
		json rankList;
		for (Who who : rank)
			rankList.push_back(who.turnFrom(Who(w)));
		msg["Rank"] = rankList;
		peer(w, msg);
		rotate(msg["Scores"]);
	}
}

std::vector<Mail> TableOpOb::popMails()
{
	std::vector<Mail> res(mMails); // copy
	mMails.clear();
	return res;
}

bool TableOpOb::gameOver() const
{
	return mEnd;
}

void TableOpOb::start()
{
	std::array<int, 4> girlIds { 0, 0, 0, 0 };
	RuleInfo rule;
	rule.roundLimit = 1;
	std::array<int, 4> points { 25000, 25000, 25000, 25000 };
	std::array<TableOperator*, 4> ops {
		&mOps[0], &mOps[1], &mOps[2], &mOps[3]
	};
	std::vector<TableObserver*> obs { this };
	Who td(0);

	mTable.reset(new Table(points, girlIds, ops, obs, rule, td));

    mTable->start();
}

void TableOpOb::action(int who, const string &actStr, const string &actArg)
{
	util::p("===TableOpOb get act", actStr, actArg);
	Action action = makeAction(actStr, actArg);
	mTable->action(Who(who), action);
}

void TableOpOb::peer(int w, const json &msg)
{
	mMails.emplace_back(w, msg.dump());
}

void TableOpOb::broad(const json &msg)
{
	const auto &s = msg.dump();
	for (int w = 0; w < 4; w++)
		mMails.emplace_back(w, s);
}



