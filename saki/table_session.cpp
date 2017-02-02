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
    return mOpOb->popMails();
}

std::vector<Mail> TableSession::Action(int who, 
                                       const std::string &actStr,
                                       const std::string &actArg) 
{
    mOpOb->action(who, actStr, actArg);
    return mOpOb->popMails();
}

std::vector<Mail> TableSession::SweepAll() 
{
    mOpOb->sweepAll();
    return mOpOb->popMails();
}

std::vector<Mail> TableSession::SweepOne(int who)
{
    mOpOb->sweepOne(who);
    return mOpOb->popMails();
}

bool TableSession::GameOver() const
{
	return mOpOb->gameOver();
}

