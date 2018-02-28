#include "table_session.h"
#include "table_opob.h"

TableSession::TableSession(int id0, int id1, int id2, int id3)
{
	std::array<int, 4> girlIds { id0, id1, id2, id3 };
    mOpOb = new TableOpOb(girlIds);
}

TableSession::~TableSession()
{
    delete mOpOb;
}

std::vector<Mail> TableSession::Start() 
{
    return mOpOb->start();
}

std::vector<Mail> TableSession::Action(int who, 
                                       const std::string &actStr,
									   int actArg,
                                       const std::string &actTile,
									   int nonce) 
{
    return mOpOb->action(who, actStr, actArg, actTile, nonce);
}

std::vector<Mail> TableSession::SweepAll() 
{
	return mOpOb->sweepAll();
}

std::vector<Mail> TableSession::SweepOne(int who)
{
    return mOpOb->sweepOne(who);
}

