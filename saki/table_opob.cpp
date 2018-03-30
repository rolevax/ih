#include "table_opob.h"

#include "libsaki/util/string_enum.h"
#include "libsaki/ai/ai.h"
#include "libsaki/util/misc.h"


using namespace saki;

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
{
    Rule rule;
    rule.roundLimit = 8;
    std::array<int, 4> points { 25000, 25000, 25000, 25000 };

	std::array<std::unique_ptr<Girl>, 4> girls;
	for (int w = 0; w < 4; w++) {
		Girl::Id id = static_cast<Girl::Id>(girlIds[w]);
		girls[w] = Girl::create(Who(w), id);
	}

    std::vector<TableObserver*> obs { &mStat, &mReplay };
    Table::InitConfig config {
        points, std::move(girls), rule, Who(0)
    };

    mServer.reset(new TableServer(std::move(config), obs, mEnv));
}

auto TableOpOb::start() -> Mails
{
    auto msgs = mServer->start();
    return mailsOfMsgs(msgs);
}

auto TableOpOb::action(int w, const string &actStr,
                       int actArg, const string &actTile,
                       int nonce) -> Mails
{
    Who who(w);

    if (actStr == "SWEEP")
        return sweepOne(w);

    if (actStr == "RESUME")
        return resume(w);

    Action action = makeAction(actStr, actArg, actTile);
    auto msgs = mServer->action(who, action, nonce);
    return mailsOfMsgs(msgs);
}

auto TableOpOb::sweepOne(int w) -> Mails
{
    Who who(w);
    const auto &choices = mServer->table().getChoices(who);
    Action act = choices.timeout();
    if (act.act() == ActCode::NOTHING)
        return Mails();

    int nonce = mServer->table().getNonce(who);
    auto msgs = mServer->action(who, act, nonce);
    return mailsOfMsgs(msgs);
}

auto TableOpOb::sweepAll() -> Mails
{
    std::array<Action, 4> actions;
    using AC = ActCode;
    for (int w = 0; w < 4; w++) {
        const auto &choices = mServer->table().getChoices(Who(w));
        actions[w] = choices.timeout();
    }

    Mails mails;

    for (int w = 0; w < 4; w++) {
        if (actions[w].act() != AC::NOTHING) {
            int nonce = mServer->table().getNonce(Who(w));
            auto subMsgs = mServer->action(Who(w), actions[w], nonce);
            auto subMails = mailsOfMsgs(subMsgs);
            mails.insert(mails.end(), subMails.begin(), subMails.end());
        }
    }

    return mails;
}

auto TableOpOb::resume(int c) -> Mails
{
    auto msgs = mServer->resume(Who(c));
    return mailsOfMsgs(msgs);
}

auto TableOpOb::mailsOfMsgs(const TableServer::Msgs &msgs) -> Mails
{
    TableOpOb::Mails mails;

    auto push = [&mails](const TableMsg &msg) {
        int to = msg.to.somebody() ? msg.to.index() : -1;
        mails.emplace_back(to, msg.content.marshal());
    };

    for (const TableMsg &msg : msgs) {
        // filter: add table-end-stat before game-over
        if (msg.content.event() == "game-over")
            push(tableEndStat());

        push(msg);
    }

    return mails;
}

auto TableOpOb::tableEndStat() -> TableMsg
{
    using json = nlohmann::json;

    json args;

    json rankList;
    for (int w = 0; w < 4; w++)
        rankList.push_back(mServer->table().getRank(Who(w)));
    args["Ranks"] = rankList;

    args["Points"] = mServer->table().getPoints();


    const auto &scores = mStat.scores();
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

    args["Replay"] = mReplay;

    return TableMsg { Who(), TableMsgContent("table-end-stat", args) };
}



