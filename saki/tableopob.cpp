#include "tableopob.h"

#include "libsaki/string_enum.h"

#include "json.hpp"

using json = nlohmann::json;

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
    const saki::TicketFolder &tifo = table.getTicketFolder(who);

	// FUCK append to msg

    json map;
    saki::Who focusWho;

    using AC = saki::ActCode;
    static const AC just[] = {
        AC::PASS, AC::SPIN_OUT, AC::RIICHI,
        AC::TSUMO, AC::RYUUKYOKU,
        AC::END_TABLE, AC::NEXT_ROUND,
        AC::DICE, AC::IRS_CLICK
    };

    for (AC code : just)
        if (tifo.can(code))
            map[saki::stringOf(code)] = true;

    json msg;
    msg["Type"] = "t-activated";
    msg["Action"] = map;
    msg["LastDiscarder"] = focusWho.nobody() ? -1 
                                             : focusWho.turnFrom(who);
    mMails.emplace_back(who.index(), msg.dump());
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

    mTable->start();
}

void TableOpOb::action(int who, int encodedAct)
{
	// decode and pass to table
}



