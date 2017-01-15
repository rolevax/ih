#include "tableopob.h"

TableOp::TableOp(TableOpOb &opOb, saki::Who self) 
	: saki::TableOperator(self)
	, mOpOb(opOb)
{
}

void TableOp::onActivated(saki::Table &table)
{
	mOpOb.onActivated(mSelf, table);
}



TableOpOb::TableOpOb()
	: mOps {
		TableOp(*this, saki::Who(0)),
		TableOp(*this, saki::Who(1)),
		TableOp(*this, saki::Who(2)),
		TableOp(*this, saki::Who(3))
	}
{
}

void TableOpOb::onActivated(saki::Who who, saki::Table &table)
{
	// FUCK append to msg
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
	saki::RuleInfo rule;
	std::array<int, 4> points { 25000, 25000, 25000, 25000 };
	std::array<saki::TableOperator*, 4> ops {
		&mOps[0], &mOps[1], &mOps[2], &mOps[3]
	};
	std::vector<saki::TableObserver*> obs { this };
	saki::Who td(0);

	mTable.reset(new saki::Table(points, girlIds, ops, obs, rule, td));
}

void TableOpOb::action(int who, int encodedAct)
{
	// decode and pass to table
}



