#ifndef SAKI_TABLESESSION_H
#define SAKI_TABLESESSION_H

#include "mail.h"

#include <vector>
#include <string>

class TableOpOb;

class TableSession
{
public:
	TableSession();
	~TableSession();
	std::vector<Mail> Start();
	std::vector<Mail> Action(int who, int encodedAct);

private:
	TableOpOb *mOpOb;
};

#endif // SAKI_TABLESESSION_H

