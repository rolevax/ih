#include "table_opob.h"

#include "json.hpp"

#include "libsaki/util/string_enum.h"
#include "libsaki/ai/ai.h"
#include "libsaki/util/misc.h"



TableOp::TableOp(TableOpOb &opOb, Who self) 
	: TableOperator(self)
	, mOpOb(opOb)
{
}

void TableOp::onActivated(Table &table)
{
	mOpOb.onActivated(mSelf, table);
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
                  


TableOpOb::TableOpOb(const std::array<int, 4> &girlIds)
	: mOps {
		TableOp(*this, Who(0)),
		TableOp(*this, Who(1)),
		TableOp(*this, Who(2)),
		TableOp(*this, Who(3))
	}
{
	Rule rule;
	rule.roundLimit = 8;
	std::array<int, 4> points { 25000, 25000, 25000, 25000 };
	std::array<TableOperator*, 4> ops {
		&mOps[0], &mOps[1], &mOps[2], &mOps[3]
	};
	std::vector<TableObserver*> obs { this, &mStat, &mReplay };
	Who td(0);

	mTable.reset(new Table(points, girlIds, ops, obs, rule, td, mEnv));

    mTable->start();
}

void TableOpOb::onActivated(Who who, Table &table)
{
    using AC = ActCode;
    using Mode = Choices::Mode;

    const auto view = table.getView(who);
    const Choices &choices = view->myChoices();

    if (table.riichiEstablished(who) && choices.spinOnly()) {
		json args;
		args["Who"] = who.index();
		system("riichi-auto", args);
		return;
    }

    json map;
    int focusWho = -1;

    switch (choices.mode()) {
    case Mode::WATCH:
        break;
    case Mode::CUT:
        activateIrsCheck(map, *view);
        break;
    case Mode::DICE:
        map[util::stringOf(AC::DICE)] = true;
        break;
    case Mode::DRAWN:
        activateDrawn(map, *view);
        break;
    case Mode::BARK:
        focusWho = view->getFocus().who().turnFrom(who);
        activateBark(map, *view);
        break;
    case Mode::END:
        if (choices.can(AC::END_TABLE))
            map[util::stringOf(AC::END_TABLE)] = true;
        if (choices.can(AC::NEXT_ROUND))
            map[util::stringOf(AC::NEXT_ROUND)] = true;
        break;
    }

    if (choices.can(AC::IRS_CLICK))
        map[util::stringOf(AC::NEXT_ROUND)] = true;

    json args;
    args["action"] = map;
    args["lastDiscarder"] = focusWho;
    args["green"] = view->myChoices().forwardAll();
	peer(who.index(), "activated", args);
}

void TableOpOb::onTableStarted(const Table &table, uint32_t seed)
{
	(void) seed;
	onPointsChanged(table);
}

void TableOpOb::onFirstDealerChoosen(Who initDealer)
{
	json args;
	for (int w = 0; w < 4; w++) {
		args["dealer"] = initDealer.turnFrom(Who(w));
		peer(w, "first-dealer-choosen", args);
	}
}

void TableOpOb::onRoundStarted(int r, int e, Who d, 
                               bool al, int dp, uint32_t s)
{
	if (r > 100) { // prevent infinite offline-spin-out loop
		mEnd = true;
		util::p("C++: round over 100");
		return;
	}

	json args;
	args["round"] = r;
	args["extra"] = e;
	args["allLast"] = al;
	args["deposit"] = dp;
	for (int w = 0; w < 4; w++) {
		args["dealer"] = d.turnFrom(Who(w));
		peer(w, "round-started", args);
	}

	args["dealer"] = d.index();
	args["seed"] = s;
	system("round-start-log", args);
}

void TableOpOb::onCleaned()
{
	broad("cleaned", json::object());
}

void TableOpOb::onDiced(const Table &table, int die1, int die2)
{
	json args;
	args["die1"] = die1;
	args["die2"] = die2;
	broad("diced", args);
}

void TableOpOb::onDealt(const Table &table)
{
	for (int w = 0; w < 4; w++) {
		const auto &init = table.getHand(Who(w)).closed().t37s13(true);
		json args;
		args["init"] = createTiles(init.range());
		peer(w, "dealt", args);
	}
}

void TableOpOb::onFlipped(const Table &table)
{
	json args;
	args["newIndic"] = createTile(table.getMount().getDrids().back());
	broad("flipped", args);
}

void TableOpOb::onDrawn(const Table &table, Who who)
{
	const T37 &in = table.getHand(who).drawn();
	for (int w = 0; w < 4; w++) {
		json args;
		args["who"] = who.turnFrom(Who(w));
		if (table.duringKan())
			args["rinshan"] = true;
		if (w == who.index())
			args["tile"] = createTile(in);
		peer(w, "drawn", args);
	}
}

void TableOpOb::onDiscarded(const Table &table, bool spin)
{
	Who discarder = table.getFocus().who();
	const T37 &out = table.getFocusTile();
	bool lay = table.lastDiscardLay();

	json args;
	args["tile"] = createTile(out, lay);
	args["spin"] = spin;
	for (int w = 0; w < 4; w++) {
		args["who"] = discarder.turnFrom(Who(w));
		peer(w, "discarded", args);
	}
}

void TableOpOb::onRiichiCalled(Who who)
{
	for (int w = 0; w < 4; w++) {
		json args;
		args["who"] = who.turnFrom(Who(w));
		peer(w, "riichi-called", args);
	}
}

void TableOpOb::onRiichiEstablished(Who who)
{
	for (int w = 0; w < 4; w++) {
		json args;
		args["who"] = who.turnFrom(Who(w));
		peer(w, "riichi-established", args);
	}
}

void TableOpOb::onBarked(const Table &table, Who who, 
                         const M37 &bark, bool spin)
{
	Who from = bark.isCpdmk() ? table.getFocus().who() : Who();

	json args;
	args["actStr"] = util::stringOf(bark.type());
	args["bark"] = createBark(bark);
	args["spin"] = spin;
	for (int w = 0; w < 4; w++) {
		args["who"] = who.turnFrom(Who(w));
		args["fromWhom"] = from.somebody() ? from.turnFrom(Who(w)) : -1;
		peer(w, "barked", args);
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
		handMap["closed"] = createTiles(hand.closed().t37s13(true).range());
		handMap["barks"] = createBarks(hand.barks());

		if (result == RR::TSUMO || result == RR::KSKP)
			handMap["pick"] = createTile(hand.drawn());
		else if (result == RR::RON || result == RR::SCHR)
			handMap["pick"] = createTile(table.getFocusTile());

		handsList.emplace_back(handMap);
	}

	for (size_t i = 0; i < forms.size(); i++) {
		const Form &form = forms[i];
		json formMap;
		formMap["spell"] = form.spell();
		formMap["charge"] = form.charge();
		formsList.emplace_back(formMap);
	}

	json args;
	args["result"] = util::stringOf(result);
	args["hands"] = handsList;
	args["forms"] = formsList;
	args["urids"] = createTiles(table.getMount().getUrids().range());
	for (int w = 0; w < 4; w++) {
		args["openers"] = json::array();
		for (Who who : openers)
			args["openers"].push_back(who.turnFrom(Who(w)));
		args["gunner"] = gunner.somebody() ? gunner.turnFrom(Who(w)) : -1;
		peer(w, "round-ended", args);
	}
}

void TableOpOb::onPointsChanged(const Table &table)
{
	json args;
	args["points"] = table.getPoints();
	for (int w = 0; w < 4; w++) {
		peer(w, "points-changed", args);
		rotate(args["points"]);
	}
}

void TableOpOb::onTableEnded(const std::array<Who, 4> &rank,
		                     const std::array<int, 4> &scores)
{
	mEnd = true;

	tableEndStat(scores);

	json args;
	args["scores"] = scores;
	for (int w = 0; w < 4; w++) {
		json rankList;
		for (Who who : rank)
			rankList.push_back(who.turnFrom(Who(w)));
		args["rank"] = rankList;
		peer(w, "table-ended", args);
		rotate(args["scores"]);
	}
}

void TableOpOb::onPoppedUp(const Table &table, Who who)
{
	json args;
	args["str"] = table.getGirl(who).popUpStr();
	peer(who.index(), "popped-up", args);
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

void TableOpOb::action(int w, const string &actStr,
                       int actArg, const string &actTile)
{
	Who who(w);

	if (actStr == "SWEEP") {
		sweepOne(w);
	} else if (actStr == "BOT") {
		Girl::Id girlId = mTable->getGirl(who).getId();
		std::unique_ptr<Ai> ai(Ai::create(who, girlId));
		ai->onActivated(*mTable);
	} else if (actStr == "RESUME") {
		resume(w);
	} else {
		Action action = makeAction(actStr, actArg, actTile, w);
		if (mTable->check(who, action)) {
			mTable->action(who, action);
		} else {
			json args;
			args["who"] = w;
			args["actStr"] = actStr;
			args["actArg"] = actArg;
			system("cannot", args);
		}
	}
}

void TableOpOb::sweepOne(int w)
{
	Who who(w);
	const auto &choices = mTable->getChoices(who);
	Action act = choices.sweep();
	if (act.act() == ActCode::NOTHING)
		return;
	mTable->action(who, act);
}

std::vector<int> TableOpOb::sweepAll()
{
	std::array<Action, 4> actions;
	using AC = ActCode;
	for (int w = 0; w < 4; w++) {
		const auto &choices = mTable->getChoices(Who(w));
		actions[w] = choices.sweep();
	}

	std::vector<int> res;
	for (int w = 0; w < 4; w++) {
		if (actions[w].act() != AC::NOTHING) {
			res.push_back(w);
			mTable->action(Who(w), actions[w]);
		}
	}
	return res;
}

void TableOpOb::resume(int c)
{
	if (mTable->beforeEast1())
		return; // no serious info needs to be provided

	json args;
	Who comer(c);

	args["whoDrawn"] = -1;
	args["barkss"] = json::array();
	args["rivers"] = json::array();
	args["riichiBars"] = json::array();
	args["dice"] = mTable->getDice();
	if (mTable->getDice() > 0) {
		for (int w = 0; w < 4; w++) {
			const Hand &hand = mTable->getHand(Who(w));
			int pers = Who(w).turnFrom(comer);
			if (hand.hasDrawn()) {
				args["whoDrawn"] = pers;
				if (w == c)
					args["drawn"] = createTile(hand.drawn());
			}
			if (w == c)
				args["myHand"] = createTiles(hand.closed().t37s13(true).range());
			args["barkss"][pers] = createBarks(hand.barks());
			args["rivers"][pers] = createTiles(mTable->getRiver(Who(w)).range());
			args["riichiBars"][pers] = mTable->riichiEstablished(Who(w));
		}
		args["drids"] = createTiles(mTable->getMount().getDrids().range());
	}

	const auto &pts = mTable->getPoints();
	args["points"] = json {
		pts[comer.index()],
		pts[comer.right().index()],
		pts[comer.cross().index()],
		pts[comer.left().index()],
	};

	args["girlIds"] = json {
		static_cast<int>(mTable->getGirl(comer).getId()),
		static_cast<int>(mTable->getGirl(comer.right()).getId()),
		static_cast<int>(mTable->getGirl(comer.cross()).getId()),
		static_cast<int>(mTable->getGirl(comer.left()).getId())
	};

	args["wallRemain"] = mTable->getMount().wallRemain();
	args["deadRemain"] = mTable->getMount().deadRemain();

	args["round"] = mTable->getRound();
	args["extraRound"] = mTable->getExtraRound();
	args["dealer"] = mTable->getDealer().turnFrom(comer);
	args["allLast"] = mTable->isAllLast();
	args["deposit"] = mTable->getDeposit();

	peer(c, "resume", args);
}

void TableOpOb::tableEndStat(const std::array<int, 4> &scores)
{
	json args;

	json rankList;
	for (int w = 0; w < 4; w++)
		rankList.push_back(mTable->getRank(Who(w)));
	args["Ranks"] = rankList;

	args["Points"] = mTable->getPoints();
	args["ATop"] = std::count_if(scores.begin(), scores.end(), 
			                     [](int s) { return s > 0; }) == 1;
	args["ALast"] = std::count_if(scores.begin(), scores.end(), 
			                      [](int s) { return s < 0; }) == 1;
	args["Round"] = mStat.roundCt();
	args["Wins"] = mStat.wins();
	args["Guns"] = mStat.guns();
	args["Barks"] = mStat.barks();
	args["Riichis"] = mStat.riichis();
	args["WinSumPoints"] = mStat.winSumPoints();
	args["GunSumPoints"] = mStat.gunSumPoints();
	args["BarkSumPoints"] = mStat.barkSumPoints();
	args["RiichiSumPoints"] = mStat.riichiSumPoints();
	args["ReadySumTurns"] = mStat.readySumTurns();
	args["Readys"] = mStat.readys();
	args["WinSumTurns"] = mStat.winSumTurns();
	args["Yakus"] = mStat.yakus();
	args["SumHans"] = mStat.sumHans();
	args["Kzeykms"] = mStat.kzeykms();

	args["Replay"] = createReplay(mReplay);

	system("table-end-stat", args);
}

void TableOpOb::peer(int w, const char *event, const json &args)
{
	json msg;
	msg["Type"] = "table";
	msg["Event"] = event;
	msg["Args"] = args;
	mMails.emplace_back(w, msg.dump());
}

void TableOpOb::broad(const char *event, const json &args)
{
	json msg;
	msg["Type"] = "table";
	msg["Event"] = event;
	msg["Args"] = args;
	const auto &s = msg.dump();
	for (int w = 0; w < 4; w++)
		mMails.emplace_back(w, s);
}

void TableOpOb::system(const char *type, const json &args)
{
	json msg = args;
	msg["Type"] = type;
	const auto &s = msg.dump();
	mMails.emplace_back(-1, s);
}

void TableOpOb::activateDrawn(json &map, const TableView &view)
{
    using AC = ActCode;

    for (AC ac : { AC::SPIN_OUT, AC::SPIN_RIICHI, AC::TSUMO, AC::RYUUKYOKU })
        if (view.myChoices().can(ac))
            map[util::stringOf(ac)] = true;

    const Choices::ModeDrawn &mode = view.myChoices().drawn();

    if (mode.swapOut)
        map[util::stringOf(AC::SWAP_OUT)] = (1 << 13) - 1;

    if (!mode.swapRiichis.empty()) {
		const auto &closed = view.myHand().closed();
        map[util::stringOf(AC::SWAP_RIICHI)] = createSwapMask(closed, mode.swapRiichis);
	}

    if (!mode.ankans.empty())
        map[util::stringOf(AC::ANKAN)] = createTileStrs(mode.ankans.range());

    if (!mode.kakans.empty()) {
		std::vector<int> kakans;
		for (int i : mode.kakans)
			kakans.push_back(i);
        map[util::stringOf(AC::KAKAN)] = kakans;
    }
}

void TableOpOb::activateBark(json &map, const TableView &view)
{
    using AC = saki::ActCode;

    std::array<AC, 7> just {
        AC::PASS,
        AC::CHII_AS_LEFT, AC::CHII_AS_MIDDLE, AC::CHII_AS_RIGHT,
        AC::PON, AC::DAIMINKAN, AC::RON
    };

    for (AC ac : just)
        if (view.myChoices().can(ac))
            map[util::stringOf(ac)] = true;
}

void TableOpOb::activateIrsCheck(json &map, const TableView &view)
{
    const Girl &girl = view.me();
    int prediceCount = girl.irsCheckCount();
    json list;
    for (int i = 0; i < prediceCount; i++)
        list.push_back(createIrsCheckRow(girl.irsCheckRow(i)));
    map[util::stringOf(saki::ActCode::IRS_CHECK)] = list;
}




