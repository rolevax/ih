#ifndef SAKI_S11N_H
#define SAKI_S11N_H

#include "libsaki/tilecount.h"
#include "libsaki/meld.h"
#include "libsaki/replay.h"

using namespace saki;

#include "json.hpp"

using json = nlohmann::json;

unsigned createSwapMask(const TileCount &closed,
                        const std::vector<T37> &choices);
std::vector<std::string> createTileStrs(const std::vector<T34> &ts);
std::string createTile(const T37 &t, bool lay = false);
json createTiles(const std::vector<T37> &ts);
json createBark(const M37 &m);
json createBarks(const std::vector<M37> &ms);

json createReplay(const Replay &replay);
json createRule(const RuleInfo &rule);
json createRound(const Replay::Round &round);
json createTrack(const Replay::Track &track);

#endif // SAKI_S11N_H



