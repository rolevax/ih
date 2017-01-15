#ifndef SAKI_TABLEOPOB_H
#define SAKI_TABLEOPOB_H

#include "mail.h"

#include "libsaki/tableoperator.h"
#include "libsaki/tableobserver.h"
#include "libsaki/table.h"

#include <memory>

class TableOpOb;

class TableOp : public saki::TableOperator
{
public:
	explicit TableOp(TableOpOb &opOb, saki::Who self);
	void onActivated(saki::Table &table) override;

private:
	TableOpOb &mOpOb;
};

class TableOpOb : public saki::TableObserver
{
public:
	TableOpOb();

	void onActivated(saki::Who who, saki::Table &table);
	virtual ~TableOpOb() = default;

	std::vector<Mail> popMails();

	void start();
	void action(int who, int encodedAct);

private:
	std::vector<Mail> mMails;
	std::unique_ptr<saki::Table> mTable;
	std::array<TableOp, 4> mOps;
};



#endif // SAKI_TABLEOPOB_H



