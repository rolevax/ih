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
	json res;
	for (const T37 &t : ts)
		res.emplace_back(createTile(t, false));
	return res;
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

	/* FUCK ,
	if (view.iCan(AC::IRS_CHECK)) {
		const Girl &girl = table.getGirl(mSelf);
		int prediceCount = girl.irsCheckCount();
		QVariantList list;
		for (int i = 0; i < prediceCount; i++)
			list << createIrsCheckRowVar(girl.irsCheckRow(i));
		map.insert(stringOf(AC::IRS_CHECK), QVariant::fromValue(list));
	}
	*/

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

std::vector<Mail> TableOpOb::popMails()
{
	std::vector<Mail> res(mMails); // copy
	mMails.clear();
	return res;
}

void TableOpOb::start()
{
	std::array<int, 4> girlIds { 0, 0, 0, 0 };
	RuleInfo rule;
	std::array<int, 4> points { 25001, 25000, 25000, 25000 };
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



