package blockingpool

// BlockingPool is a generic, channel-based object pool that provides blocking
// semantics for both acquiring and returning objects.
//
// The pool has a fixed capacity, specified at creation time. It is for
// scenarios where you want to limit the number of concurrently allocated
// resources and enforce strict back-pressure:
//
//   - Get() blocks until an object is available in the pool (or the caller’s
//     context is canceled if used with select).
//   - Put() blocks until there is space in the pool (i.e., the number of
//     outstanding objects is below the capacity).
//
// Important characteristics:
//   - Get() will block indefinitely if the pool is empty or until an new item
//     is .Put() into the pool.
//   - Put() will block indefinitely if the pool is at full capacity or until
//     an item is .Get() from the pool.
type BlockingPool[T any] struct {
	pool chan T
}

// NewBlockingPool creates a new BlockingPool with the specified capacity.
//
// The capacity determines the maximum number of objects that can be "checked
// out" simultaneously (i.e., the maximum number of outstanding Get() calls
// without corresponding Put() calls).
func NewBlockingPool[T any](capacity int) BlockingPool[T] {
	return BlockingPool[T]{pool: make(chan T, capacity)}
}

// Get acquires an object from the pool, blocking until one is available.
//
// The returned value is whatever was previously Put into the pool. If the pool
// is empty, .Get() will block indefinitely until another goroutine calls
// .Put().
//
// It is the caller's responsibility to eventually call .Put() with the
// returned object (or a replacement) to release it back to the pool.
func (p *BlockingPool[T]) Get() T { return <-p.pool }

// Put returns an object to the pool, blocking until there is space available.
//
// If the pool is already at full capacity, .Put() will block until another
// goroutine calls .Get().
//
// After a successful Put(), the object becomes available for .Get() calls.
func (p *BlockingPool[T]) Put(obj T) { p.pool <- obj }
