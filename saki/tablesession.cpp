#include "tablesession.h"
#include "tableopob.h"

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

std::vector<Mail> TableSession::Action(int who, int encodedAct) 
{
	mOpOb->action(who, encodedAct);
	return mOpOb->popMails();
}

