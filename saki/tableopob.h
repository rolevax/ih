#ifndef SAKI_TABLEOPOB_H
#define SAKI_TABLEOPOB_H

#include "mail.h"

#include "json.hpp"

#include "libsaki/tableoperator.h"
#include "libsaki/tableobserver.h"
#include "libsaki/table.h"

#include <memory>

class TableOpOb;

using namespace saki;
using string = std::string;
using json = nlohmann::json;

class TableOp : public TableOperator
{
public:
	explicit TableOp(TableOpOb &opOb, Who self);

	void onActivated(Table &table) override;

private:
	TableOpOb &mOpOb;
};

class TableOpOb : public TableObserver
{
public:
	TableOpOb();
	virtual ~TableOpOb() = default;

	void onActivated(Who who, Table &table);
	void onTableStarted(const Table &table, uint32_t seed) override;
	void onFirstDealerChoosen(Who initDealer) override;
    void onRoundStarted(int r, int e, Who d, 
                        bool al, int dp, uint32_t s) override;
	void onCleaned() override;
	void onDiced(const Table &table, int die1, int die2) override;
	void onDealt(const Table &table) override;
	void onFlipped(const Table &talbe) override;
	void onDrawn(const Table &table, Who who) override;
	void onPointsChanged(const Table &table) override;

	std::vector<Mail> popMails();


	void start();
	void action(int who, const string &actStr, const string &actArg);

private:
	void peer(int w, const json &msg);
	void broad(const json &msg);

private:
	std::vector<Mail> mMails;
	std::unique_ptr<Table> mTable;
	std::array<TableOp, 4> mOps;
};



#endif // SAKI_TABLEOPOB_H



