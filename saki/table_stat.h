#ifndef SAKI_TABLE_STAT_H
#define SAKI_TABLE_STAT_H

#include "libsaki/tableobserver.h"



using namespace saki;



class TableStat : public TableObserver
{
public:
	TableStat() = default;
	virtual ~TableStat() = default;

	void onRoundEnded(const Table &table, RoundResult result,
                      const std::vector<Who> &openers, Who gunner,
					  const std::vector<Form> &forms) override;
	void onTableEnded(const std::array<Who, 4> &rank,
			          const std::array<int, 4> &scores) override;

private:

private:
};



#endif // SAKI_TABLE_STAT_H



