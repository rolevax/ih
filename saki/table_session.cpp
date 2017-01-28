#include "table_session.h"
#include "table_opob.h"

TableSession::TableSession()
{
    mOpOb = new TableOpOb();
}

TableSession::~TableSession()
{
    delete mOpOb;
}

std::vector<Mail> TableSession::Start() 
{
    mOpOb->start();
    return mOpOb->popMails();
}

std::vector<Mail> TableSession::Action(int who, 
                                       const std::string &actStr,
                                       const std::string &actArg) 
{
    mOpOb->action(who, actStr, actArg);
    return mOpOb->popMails();
}

std::vector<Mail> TableSession::Sweep() 
{
    mOpOb->sweep();
    return mOpOb->popMails();
}

bool TableSession::GameOver() const
{
	return mOpOb->gameOver();
}

