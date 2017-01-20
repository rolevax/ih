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

	using str = std::string;

	std::vector<Mail> Start();
	std::vector<Mail> Action(int who, const str &actStr, const str &actArg);
	bool GameOver() const;

private:
	TableOpOb *mOpOb;
};

#endif // SAKI_TABLESESSION_H

