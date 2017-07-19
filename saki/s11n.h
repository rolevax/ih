#ifndef SAKI_S11N_H
#define SAKI_S11N_H

#include "libsaki/tile_count.h"
#include "libsaki/meld.h"
#include "libsaki/replay.h"

using namespace saki;

#include "json.hpp"

using json = nlohmann::json;

unsigned createSwapMask(const TileCount &closed,
                        const util::Stactor<T37, 13> &choices);
std::vector<std::string> createTileStrs(const util::Range<T34> &ts);
std::string createTile(const T37 &t, bool lay = false);
json createTiles(const std::vector<T37> &ts);
json createTiles(const util::Range<T37> &ts);
json createBark(const M37 &m);
json createBarks(const util::Stactor<M37, 4> &ms);
json createIrsCheckRow(const IrsCheckRow &row);

json createReplay(const Replay &replay);
json createRule(const RuleInfo &rule);
json createRound(const Replay::Round &round);
json createTrack(const Replay::Track &track);

Action makeAction(const std::string &actStr, int actArg,
                  const std::string &actTile, int who);

#endif // SAKI_S11N_H



