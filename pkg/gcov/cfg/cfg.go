package cfg

import gcovraw "github.com/yhlooo/gcovgo/pkg/gcov/raw"

// BuildCFG 构建控制流图
func BuildCFG(n int, blockArcs []gcovraw.RecordArcs, counts []uint64) (CFG, error) {
	cfg := make(CFG, n)
	for i := range cfg {
		cfg[i].no = uint32(i)
	}

	var setCountArcs []*Arc
	for _, b := range blockArcs {
		for _, arc := range b.Arcs {
			outArc := NewArc(cfg.Get(b.BlockNo), cfg.Get(arc.DestBlock))
			if !arc.Flags.OnTree() {
				setCountArcs = append(setCountArcs, outArc)
			}
		}
	}

	for _, arc := range setCountArcs {
		if len(counts) == 0 {
			break
		}
		arc.SetCount(counts[0])
		counts = counts[1:]
	}

	return cfg, nil
}

// CFG 控制流图
type CFG []Block

// Get 获取指定编号的块
func (cfg CFG) Get(i uint32) *Block {
	if i >= uint32(len(cfg)) {
		return nil
	}
	return &cfg[i]
}

// Block 块
type Block struct {
	// 块编号
	no uint32
	// 块执行次数
	count uint64
	// 执行次数是否已确定
	resolved bool

	// 入边
	in []*Arc
	// 出边
	out []*Arc
}

// No 返回块编号
func (blk *Block) No() uint32 {
	return blk.no
}

// Count 返回执行次数
func (blk *Block) Count() uint64 {
	return blk.count
}

// Resolved 返回执行次数是否已经确定
func (blk *Block) Resolved() bool {
	return blk.resolved
}

// In 返回入边
func (blk *Block) In() []*Arc {
	if blk.in == nil {
		return nil
	}
	ret := make([]*Arc, len(blk.in))
	copy(ret, blk.in)
	return ret
}

// Out 返回出边
func (blk *Block) Out() []*Arc {
	if blk.out == nil {
		return nil
	}
	ret := make([]*Arc, len(blk.out))
	copy(ret, blk.out)
	return ret
}

// Resolve 根据入边和出边执行次数推断块执行次数
func (blk *Block) Resolve() bool {
	//if blk.resolved {
	//	return blk.resolved
	//}

	// 通过入边推断执行次数
	count := uint64(0)
	allArcResolved := blk.in != nil
	for _, arc := range blk.in {
		if !arc.Resolved() {
			allArcResolved = false
			break
		}
		count += arc.Count()
	}

	if !allArcResolved {
		// 通过出边推断执行次数
		count = 0
		allArcResolved = blk.out != nil
		for _, arc := range blk.out {
			if !arc.Resolved() {
				allArcResolved = false
				break
			}
			count += arc.Count()
		}
	}

	// 推断不出来
	if !allArcResolved {
		return false
	}

	blk.count = count
	blk.resolved = true

	// 尝试推断入边和出边
	for _, arc := range blk.in {
		arc.Resolve()
	}
	for _, arc := range blk.out {
		arc.Resolve()
	}

	return blk.resolved
}

// NewArc 创建边
func NewArc(src, dst *Block) *Arc {
	arc := &Arc{
		src: src,
		dst: dst,
	}
	src.out = append(src.out, arc)
	dst.in = append(dst.in, arc)
	return arc
}

// Arc 边
type Arc struct {
	// 边执行次数
	count uint64
	// 执行次数是否已确定
	resolved bool

	// 源块
	src *Block
	// 目标块
	dst *Block
}

// Count 返回执行次数
func (arc *Arc) Count() uint64 {
	return arc.count
}

// Resolved 返回执行次数是否已确定
func (arc *Arc) Resolved() bool {
	return arc.resolved
}

// Source 返回源块
func (arc *Arc) Source() *Block {
	return arc.src
}

// Destination 返回目标块
func (arc *Arc) Destination() *Block {
	return arc.dst
}

// SetCount 设置执行次数
func (arc *Arc) SetCount(count uint64) {
	if arc.resolved {
		return
	}

	arc.count = count
	arc.resolved = true

	// 尝试推断源块和目标块
	if arc.src != nil {
		arc.src.Resolve()
	}
	if arc.dst != nil {
		arc.dst.Resolve()
	}
}

// Resolve 根据源和目标块执行次数推断边执行次数
func (arc *Arc) Resolve() bool {
	if arc.resolved {
		return arc.resolved
	}

	if arc.src != nil {
		if ok := arc.resolveFromNeighbor(arc.src, arc.src.Out()); ok {
			return true
		}
	}
	if arc.dst != nil {
		if ok := arc.resolveFromNeighbor(arc.dst, arc.dst.In()); ok {
			return true
		}
	}

	return false
}

// resolveFromNeighbor 根据临近块和边推断当前边执行次数
func (arc *Arc) resolveFromNeighbor(blk *Block, arcs []*Arc) bool {
	if !blk.Resolved() {
		return false
	}

	// 唯一出或入边
	if len(arcs) == 1 {
		arc.SetCount(blk.Count())
		return true
	}

	// 检查是否最后一个未确定出或入边
	count := blk.Count()
	for _, a := range arcs {
		if a == arc {
			// 跳过自己
			continue
		}
		if !a.Resolved() {
			return false
		}
		count -= a.Count()
	}
	arc.SetCount(count)
	return true
}
