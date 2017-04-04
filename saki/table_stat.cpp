#include "table_opob.h"



TableStat::TableStat()
{
	mWins.fill(0);
	mGuns.fill(0);
	mBarks.fill(0);
	mRiichis.fill(0);
}

void TableStat::onRoundEnded(const Table &table, RoundResult result,
		                     const std::vector<Who> &openers, Who gunner,
		                     const std::vector<Form> &forms)
{
	using RR = RoundResult;

	mRoundCt++;

	if (result == RR::TSUMO || result == RR::RON || result == RR::SCHR)
		for (Who who : openers)
			mWins[who.index()]++;

	if (result == RR::RON || result == RR::SCHR)
		mGuns[gunner.index()]++;

	for (int w = 0; w < 4; w++) {
		if (!table.getHand(Who(w)).isMenzen()) {
			mBarks[w]++;
			// TODO expect of bark
		} else if (table.riichiEstablished(Who(w))) {
			mRiichis[w]++;
			// TODO expect of riichi
		}
	}
	// FUCK
}

void TableStat::onTableEnded(const std::array<Who, 4> &rank,
		                     const std::array<int, 4> &scores)
{
	// FUCK
}


int TableStat::roundCt() const
{
	return mRoundCt;
}

const std::array<int, 4> &TableStat::wins() const
{
	return mWins;
}

const std::array<int, 4> &TableStat::guns() const
{
	return mGuns;
}

const std::array<int, 4> &TableStat::barks() const
{
	return mBarks;
}

const std::array<int, 4> &TableStat::riichis() const
{
	return mRiichis;
}

