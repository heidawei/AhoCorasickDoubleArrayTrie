package ac

type State struct {
	// 模式串的长度，也是这个状态的深度
	depth   int
	code    int
	// fail函数，如果没有匹配到，则跳转到此状态
	failure_ *State

	// 只要这个状态可达，则记录模式串
	emits   *TreeSet

	success map[rune]*State
	// 在base中的索引
    index int
}

func NewState(depth int) *State {
	return &State{depth: depth, success: make(map[rune]*State)}
}


func (s *State)getDepth()int {
	return s.depth
}

// 添加一个匹配的模式串(这个状态对应着这个模式串)
func (s *State)addEmit(keyword int) {
	if s.emits == nil {
		s.emits = NewTreeSet()
	}
	s.emits.Add(keyword)
}

func (s *State) getLargestValueId() int {
	if s.emits == nil || s.emits.Size() == 0 {
		return 0
	}
	return s.emits.Max()
}

func (s *State)addEmits(emits []int) {
	for _, e := range emits {
		s.addEmit(e)
	}
}

func (s *State) isAcceptable() bool {
	return s.depth > 0 && s.emits != nil
}

func (s *State) failure() *State {
	return s.failure_
}

func (s *State) setFailure(failState *State, fail []int) {
	s.failure_ = failState
	fail[s.index] = failState.index
}

func (s *State) nextState_(c rune, ignoreRootState bool) *State {
	nextState := s.success[c]
	if !ignoreRootState && nextState == nil && s.depth == 0 {
		nextState = s
	}
	return nextState
}

func (s *State) nextState(c rune) *State {
	return s.nextState_(c, false)
}

func (s *State) nextStateIgnoreRootState(c rune) *State {
	return s.nextState_(c, true)
}

func (s *State) addState(c rune) *State {
	nextState := s.nextStateIgnoreRootState(c)
	if nextState == nil {
		nextState = NewState(s.depth + 1)
		s.success[c]=nextState
	}
	return nextState
}

func (s *State)getStates() []*State {
	var states []*State
	for _, e := range s.success {
		states = append(states, e)
	}
	return states
}

func (s *State)getTransitions()[]rune {
	var ts []rune
	for k, _ := range s.success {
		ts = append(ts, k)
	}
	return ts
}

func (s *State) getSuccess() map[rune]*State{
	return s.success
}

func (s *State) getIndex() int {
	return s.index
}

func (s *State) setIndex(index int) {
	s.index = index
}

type ListState struct {
	size_   int
	states  []*State
}

func NewListState() *ListState {
	return &ListState{size_: 0}
}

func (l *ListState) size() int {
	return l.size_
}

// TODO check index > size
func (l *ListState) get(index int) *State {
	if index < 0 {
		return nil
	}
	return l.states[index]
}

func (l *ListState) add(s *State) {
	l.states = append(l.states, s)
	l.size_++
}

