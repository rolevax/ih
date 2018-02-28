#ifndef SAKI_TABLESESSION_H
#define SAKI_TABLESESSION_H

#include "mail.h"

#include <vector>
#include <string>

class TableOpOb;

class TableSession
{
public:
	TableSession(int id0, int id1, int id2, int id3);
	~TableSession();

	std::vector<Mail> Start();
	std::vector<Mail> Action(int who, const std::string &actStr, 
                             int actArg, const std::string &actTile,
							 int nonce);
	std::vector<Mail> SweepAll();
	std::vector<Mail> SweepOne(int who);

private:
	TableOpOb *mOpOb;
};

#endif // SAKI_TABLESESSION_H

