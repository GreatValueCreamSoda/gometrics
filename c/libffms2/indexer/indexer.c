#include <indexer.h>

extern int goIndexCallback(int64_t Current, int64_t Total, void *ICPrivate);

int cIndexingCallback(int64_t Current, int64_t Total, void *ICPrivate) {
    return goIndexCallback(Current, Total, ICPrivate);
}