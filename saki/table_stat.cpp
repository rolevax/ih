#include "table_opob.h"



/*
TableStat::TableStat()
{
    mWins.fill(0);
    mGuns.fill(0);
    mBarks.fill(0);
    mRiichis.fill(0);
    mWinSumPoints.fill(0);
    mGunSumPoints.fill(0);
    mBarkSumPoints.fill(0);
    mRiichiSumPoints.fill(0);
    mReadys.fill(0);
    mReadySumTurns.fill(0);
	mWinSumTurns.fill(0);
	mKzeykms.fill(0);
}

void TableStat::onRoundStarted(int r, int e, Who d,
                               bool al, int dp, uint32_t s)
{
    (void) r; (void) e; (void) d; (void) al; (void) dp; (void) s;
    mReadyMarkeds.fill(false);
}

void TableStat::onDealt(const Table &table)
{
    for (int w = 0; w < 4; w++) {
        const Hand &hand = table.getHand(Who(w));
        if (!mReadyMarkeds[w] && hand.ready()) {
			mReadySumTurns[w] += 0; // just mean this
			mReadys[w]++;
			mReadyMarkeds[w] = true;
        }
    }
}

void TableStat::onDrawn(const Table &table, Who who)
{
    const Hand &hand = table.getHand(who);
    if (!mReadyMarkeds[who.index()] && hand.ready()) {
		mReadyMarkeds[who.index()] = true;
		mReadys[who.index()]++;
		mReadySumTurns[who.index()] += table.getRiver(who).size() + 1;
    }
}

void TableStat::onBarked(const Table &table, Who who,
                         const M37 &bark, bool spin)
{
    (void) bark; (void) spin;
    const Hand &hand = table.getHand(who);
    if (!mReadyMarkeds[who.index()] && hand.ready()) {
		mReadyMarkeds[who.index()] = true;
		mReadys[who.index()]++;
		mReadySumTurns[who.index()] += table.getRiver(who).size() + 1;
    }
}

void TableStat::onRoundEnded(const Table &table, RoundResult result,
                             const std::vector<Who> &openers, Who gunner,
                             const std::vector<Form> &forms)
{
    using RR = RoundResult;

    std::array<int, 4> deltas;
    deltas.fill(0);

    mRoundCt++;

    // excluding SCHR case, counting only real-gain cases
    if (result == RR::TSUMO || result == RR::RON) {
        for (size_t i = 0; i < forms.size(); i++) {
            Who who(openers[i]);
            const Form &form = forms[i];
            mWins[who.index()]++;
            int gain = form.gain();
            mWinSumPoints[who.index()] += gain;
            deltas[who.index()] += gain;
			int turn = table.getRiver(who).size();
			if (result == RR::TSUMO)
				turn++;
			mWinSumTurns[who.index()] += turn;

			std::vector<const char*> keys = form.keys();
			int han = form.han();
			std::map<const char*, int> &yaku = mYakus[who.index()];
			std::map<const char*, int> &sumHan = mSumHans[who.index()];
			for (const char *key : keys) {
				if (yaku.find(key) == yaku.end()) {
					yaku[key] = 1;
					if (!form.isPrototypalYakuman())
						sumHan[key] = han;
				} else {
					yaku[key]++;
					if (!form.isPrototypalYakuman())
						sumHan[key] += han;
				}
			}

			auto update = [&](const char *key, int ct) {
				if (ct == 0)
					return;
				if (yaku.find(key) == yaku.end())
					yaku[key] = 0;
				yaku[key] += ct;
			};

			if (!form.isPrototypalYakuman() && han >= 13)
				mKzeykms[who.index()]++;

			if (form.dora() > 0) {
				const auto &drids = table.getMount().getDrids();
				int hyou = drids[0] % table.getHand(who);
				if (result == RR::RON)
					hyou += drids[0] % table.getFocusTile();
				update("dora", hyou);
				int kan = 0;
				for (size_t i = 1; i < drids.size(); i++) {
					kan += drids[i] % table.getHand(who);
					if (result == RR::RON)
						kan += drids[i] % table.getFocusTile();
				}
				update("kandora", kan);
			}

			if (form.uradora() > 0) {
				const auto &urids = table.getMount().getUrids();
				int hyou = urids[0] % table.getHand(who);
				if (result == RR::RON)
					hyou += urids[0] % table.getFocusTile();
				update("uradora", hyou);
				int kan = 0;
				for (size_t i = 1; i < urids.size(); i++) {
					kan += urids[i] % table.getHand(who);
					if (result == RR::RON)
						kan += urids[i] % table.getFocusTile();
				}
				update("kanuradora", kan);
			}

			if (form.akadora() > 0) {
				update("akadora", form.akadora());
			}
        }
    }

    // excluding SCHR case, counting only real-loss cases
    if (result == RR::RON) {
        mGuns[gunner.index()]++;
        int sumLoss = 0;
        for (const Form &form : forms)
            sumLoss += form.loss(false);
        mGunSumPoints[gunner.index()] += sumLoss;
        deltas[gunner.index()] -= sumLoss;
    }

    if (result == RR::HP) {
        std::array<bool, 4> tenpai { false, false, false, false };
        for (Who who : openers)
            tenpai[who.index()] = true;

        int ct = openers.size();
        if (ct % 4 != 0) {
            for (int w = 0; w < 4; w++)
                deltas[w] += tenpai[w] ? (3000 / ct) : -(3000 / (4 - ct));
        }
    } else if (result == RR::NGSMG) {
        for (Who who : openers) {
            Who dealer = table.getDealer();
            if (who == dealer) {
                for (int l = 0; l < 4; l++)
                    deltas[l] += l == who.index() ? 12000 : -4000;
            } else {
                for (int l = 0; l < 4; l++) {
                    if (l == who.index())
                        deltas[l] += 8000;
                    else
                        deltas[l] -= (Who(l) == dealer ? 4000 : 2000);
                }
            }
        }
    }

    for (int w = 0; w < 4; w++) {
        if (!table.getHand(Who(w)).isMenzen()) {
            mBarks[w]++;
            mBarkSumPoints[w] += deltas[w];
        } else if (table.riichiEstablished(Who(w))) {
            mRiichis[w]++;
            mRiichiSumPoints[w] += deltas[w];
        }
    }
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

const std::array<int, 4> &TableStat::winSumPoints() const
{
    return mWinSumPoints;
}

const std::array<int, 4> &TableStat::gunSumPoints() const
{
    return mGunSumPoints;
}

const std::array<int, 4> &TableStat::barkSumPoints() const
{
    return mBarkSumPoints;
}

const std::array<int, 4> &TableStat::riichiSumPoints() const
{
    return mRiichiSumPoints;
}

const std::array<int, 4> &TableStat::readySumTurns() const
{
	return mReadySumTurns;
}

const std::array<int, 4> &TableStat::readys() const
{
	return mReadys;
}

const std::array<int, 4> &TableStat::winSumTurns() const
{
	return mWinSumTurns;
}

const std::array<std::map<const char*, int>, 4> &TableStat::yakus() const
{
	return mYakus;
}

const std::array<std::map<const char*, int>, 4> &TableStat::sumHans() const
{
	return mSumHans;
}

const std::array<int, 4> &TableStat::kzeykms() const
{
	return mKzeykms;
}
*/

