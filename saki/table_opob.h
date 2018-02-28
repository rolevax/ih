#ifndef SAKI_TABLEOPOB_H
#define SAKI_TABLEOPOB_H

#include "mail.h"
//FUCK
//#include "table_stat.h"

#include "libsaki/app/table_server.h"
#include "libsaki/app/replay.h"
#include "libsaki/table/table_env_stub.h" // temp

#include <memory>

class TableOpOb;

using string = std::string;

class TableOpOb
{
public:
	TableOpOb(const std::array<int, 4> &girlIds);
	~TableOpOb() = default;

	using Mails = std::vector<Mail>;

	Mails start();
	Mails action(int who, const string &actStr,
                int argArg, const string &actTile,
				int nonce);
	Mails sweepAll();
	Mails sweepOne(int who);
	Mails resume(int who);

private:
	void tableEndStat(const std::array<int, 4> &scores);

private:
	//TableStat mStat;
	saki::TableEnvStub mEnv; // temp
	saki::Replay mReplay;
	std::vector<Mail> mMails;
	std::unique_ptr<saki::TableServer> mServer;
};



#endif // SAKI_TABLEOPOB_H



