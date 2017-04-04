#ifndef SAKI_TABLE_STAT_H
#define SAKI_TABLE_STAT_H

#include "libsaki/tableobserver.h"



using namespace saki;



class TableStat : public TableObserver
{
public:
	TableStat();
	virtual ~TableStat() = default;

	void onRoundEnded(const Table &table, RoundResult result,
                      const std::vector<Who> &openers, Who gunner,
					  const std::vector<Form> &forms) override;
	void onTableEnded(const std::array<Who, 4> &rank,
			          const std::array<int, 4> &scores) override;

	int roundCt() const;
	const std::array<int, 4> &wins() const;
	const std::array<int, 4> &guns() const;
	const std::array<int, 4> &barks() const;
	const std::array<int, 4> &riichis() const;

private:

private:
	int mRoundCt = 0;
	std::array<int, 4> mWins;
	std::array<int, 4> mGuns;
	std::array<int, 4> mBarks;
	std::array<int, 4> mRiichis;
};



#endif // SAKI_TABLE_STAT_H



