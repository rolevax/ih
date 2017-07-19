#include "s11n.h"

#include "libsaki/string_enum.h"

#include <bitset>

unsigned createSwapMask(const TileCount &closed,
                        const util::Stactor<T37, 13> &choices)
{
	// assume 'choices' is 34-sorted
	std::bitset<13> mask;
	int i = 0;

	auto it = choices.begin();
	for (const T37 &t : tiles37::ORDER37) {
		if (it == choices.end())
			break;
		int ct = closed.ct(t);
		if (ct > 0) {
			bool val = t.looksSame(*it);
			while (ct --> 0)
				mask[i++] = val;
			it += val; // consume choice if matched
		}
	}

	return mask.to_ulong();
}

std::vector<std::string> createTileStrs(const util::Range<T34> &ts)
{
	std::vector<std::string> res;
	for (T34 t : ts)
		res.emplace_back(t.str());
	return res;
}

std::string createTile(const T37 &t, bool lay)
{
	std::string res(t.str());
	if (lay)
		res += '_';
	return res;
}

json createTiles(const std::vector<T37> &ts)
{
	json res = json::array();
	for (const T37 &t : ts)
		res.emplace_back(createTile(t, false));
	return res;
}

json createTiles(const util::Range<T37> &ts)
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

json createBarks(const util::Stactor<M37, 4> &ms)
{
	json list = json::array();
	for (const M37 &m : ms)
		list.emplace_back(createBark(m));
	return list;
}

json createIrsCheckRow(const IrsCheckRow &row)
{
    json map;

    map["modelMono"] = row.mono;
    map["modelIndent"] = row.indent;
    map["modelText"] = row.name;
    map["modelAble"] = row.able;
    map["modelOn"] = row.on;

    return map;
}

json createReplay(const Replay &replay)
{
	json root;

    root["version"] = 3;

	json girls;
    for (saki::Girl::Id v : replay.girls)
		girls.push_back(static_cast<int>(v));
    root["girls"] = girls;

    root["initPoints"] = replay.initPoints;
    root["rule"] = createRule(replay.rule);
    root["seed"] = std::to_string(replay.seed);

    json rounds;
    for (const saki::Replay::Round &round : replay.rounds)
        rounds.push_back(createRound(round));
    root["rounds"] = rounds;

    return root;

}

json createRule(const RuleInfo &rule)
{
	json obj;

    obj["fly"] = rule.fly;
    obj["headJump"] = rule.headJump;
    obj["nagashimangan"] = rule.nagashimangan;
    obj["ippatsu"] = rule.ippatsu;
    obj["uradora"] = rule.uradora;
    obj["kandora"] = rule.kandora;
    obj["akadora"] = static_cast<int>(rule.akadora);
    obj["daiminkanPao"] = rule.daiminkanPao;
    obj["hill"] = rule.hill;
    obj["returnLevel"] = rule.returnLevel;
    obj["roundLimit"] = rule.roundLimit;

    return obj;
}

json createRound(const Replay::Round &round)
{
	json obj;

    obj["round"] = round.round;
    obj["extraRound"] = round.extraRound;
    obj["dealer"] = round.dealer.index();
    obj["allLast"] = round.allLast;
    obj["deposit"] = round.deposit;
    obj["state"] = std::to_string(round.state);
    obj["die1"] = round.die1;
    obj["die2"] = round.die2;
    obj["result"] = stringOf(round.result);
    obj["resultPoints"] = round.resultPoints;

    obj["spells"] = round.spells;
    obj["charges"] = round.charges;

    json drids;
    for (const saki::T37 &t : round.drids)
        drids.push_back(t.str());
    obj["drids"] = drids;

    json urids;
    for (const saki::T37 &t : round.urids)
        urids.push_back(t.str());
    obj["urids"] = urids;

    obj["tracks"] = json {
        createTrack(round.tracks[0]), createTrack(round.tracks[1]),
        createTrack(round.tracks[2]), createTrack(round.tracks[3])
    };

    return obj;
}

json createTrack(const Replay::Track &track)
{
	using str = std::string;

	auto inJson = [](Replay::InAct inAct) -> str {
        using In = Replay::In;
        switch (inAct.act) {
        case In::DRAW:
            return inAct.t37.str();
        case In::CHII_AS_LEFT: // 'b' means 'begin'
            return str("b") + std::to_string(inAct.showAka5);
        case In::CHII_AS_MIDDLE: // 'm' means 'middle'
            return str("m") + std::to_string(inAct.showAka5);
        case In::CHII_AS_RIGHT: // 'e' means 'end'
            return str("e") + std::to_string(inAct.showAka5);
        case In::PON:
            return str("p") + std::to_string(inAct.showAka5);
        case In::DAIMINKAN:
            return "d";
        case In::RON:
            return "r";
        case In::SKIP_IN:
            return "--";
        default:
            return "err";
        }
    };

    auto outJson = [](saki::Replay::OutAct outAct) -> str {
        using Out = Replay::Out;
        switch (outAct.act) {
        case Out::ADVANCE:
            return outAct.t37.str();
        case Out::SPIN:
            return "->";
        case Out::RIICHI_ADVANCE:
            return str("!") + outAct.t37.str();
        case Out::RIICHI_SPIN:
            return "!->";
        case Out::ANKAN:
            return str("a") + outAct.t37.str();
        case Out::KAKAN:
            return str("k") + outAct.t37.str();
        case Out::RYUUKYOKU:
            return "~";
        case Out::TSUMO:
            return "t";
        case Out::SKIP_OUT:
            return "--";
        default:
            return "err";
        }
    };

    json obj;

    json initArr;
    for (const saki::T37 &t : track.init)
        initArr.push_back(t.str());
    obj["init"] = initArr;

    json inArr;
    for (const saki::Replay::InAct &inAct : track.in)
        inArr.push_back(inJson(inAct));
    obj["in"] = inArr;

    json outArr;
    for (const saki::Replay::OutAct &outAct : track.out)
        outArr.push_back(outJson(outAct));
    obj["out"] = outArr;

    return obj;
}



Action makeAction(const std::string &actStr, int actArg,
                  const std::string &actTile, int who)
{
    using AC = saki::ActCode;
	AC act = actCodeOf(actStr.c_str());

	switch (act) {
		case AC::SWAP_OUT:
		case AC::SWAP_RIICHI:
			return Action(act, T37(actTile.c_str()));
		case AC::ANKAN:
			return Action(act, T34(actTile.c_str()));
		case AC::CHII_AS_LEFT:
		case AC::CHII_AS_MIDDLE:
		case AC::CHII_AS_RIGHT:
		case AC::PON:
			return Action(act, actArg, T37(actTile.c_str()));
		case AC::KAKAN:
			return Action(act, static_cast<int>(actArg));
		case AC::IRS_CHECK:
			return Action(act, static_cast<unsigned>(actArg));
		default:
			return Action(act);
	}
}

