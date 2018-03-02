#ifndef SAKI_TABLE_STAT_H
#define SAKI_TABLE_STAT_H

#include "libsaki/table/table_observer.h"

#include <map>



using namespace saki;



class TableStat : public TableObserverDispatched
{
public:
    TableStat();
    virtual ~TableStat() = default;

    void onTableEvent(const Table &table, const TE::RoundStarted &event) override;
    void onTableEvent(const Table &table, const TE::Dealt &event) override;
    void onTableEvent(const Table &table, const TE::Drawn &event) override;
    void onTableEvent(const Table &table, const TE::Barked &event) override;
    void onTableEvent(const Table &table, const TE::RoundEnded &event) override;
    void onTableEvent(const Table &table, const TE::TableEnded &event) override;

    int roundCt() const;
    const std::array<int, 4> &wins() const;
    const std::array<int, 4> &guns() const;
    const std::array<int, 4> &barks() const;
    const std::array<int, 4> &riichis() const;
    const std::array<int, 4> &winSumPoints() const;
    const std::array<int, 4> &gunSumPoints() const;
    const std::array<int, 4> &barkSumPoints() const;
    const std::array<int, 4> &riichiSumPoints() const;
    const std::array<int, 4> &readySumTurns() const;
    const std::array<int, 4> &readys() const;
    const std::array<int, 4> &winSumTurns() const;
    const std::array<std::map<const char*, int>, 4> &yakus() const;
    const std::array<std::map<const char*, int>, 4> &sumHans() const;
    const std::array<int, 4> &kzeykms() const;
    const std::array<int, 4> &scores() const;

private:
    int mRoundCt = 0;
    std::array<int, 4> mWins;
    std::array<int, 4> mGuns;
    std::array<int, 4> mBarks;
    std::array<int, 4> mRiichis;
    std::array<int, 4> mWinSumPoints;
    std::array<int, 4> mGunSumPoints;
    std::array<int, 4> mBarkSumPoints;
    std::array<int, 4> mRiichiSumPoints;
    std::array<bool, 4> mReadyMarkeds; // temp one-round use
    std::array<int, 4> mReadys;
    std::array<int, 4> mReadySumTurns;
    std::array<int, 4> mWinSumTurns;
    std::array<std::map<const char*, int>, 4> mYakus;
    std::array<std::map<const char*, int>, 4> mSumHans;
    std::array<int, 4> mKzeykms;
    std::array<int, 4> mScores;
};



#endif // SAKI_TABLE_STAT_H



