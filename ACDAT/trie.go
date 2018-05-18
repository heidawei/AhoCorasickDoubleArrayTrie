package ACDAT

import (
	"sort"
	"fmt"
)

type Word struct {
	runes  []rune
}

func NewWord(word string) *Word {
	return &Word{runes: []rune(word)}
}

func (w *Word)GetWord() string {
	return string(w.runes)
}

func (w *Word)GetRune(index int) rune {
	if index < 0 || index >= len(w.runes) {
		panic("invalid index")
	}
	return w.runes[index]
}

func (w *Word)GetRunes() []rune {
	return w.runes
}

func (w *Word) Size() int {
	return len(w.runes)
}

func (w *Word)String() string {
	return string(w.runes)
}

type WordCodeDict struct {
	nextCode int
	dict     map[rune]int
}

func NewWordCodeDict(words []*Word) *WordCodeDict {
	var _words []rune
	dict := make(map[rune]int)
	for _, word := range words {
		for _, r := range word.GetRunes() {
			dict[r] = 0
		}
	}
	for r, _ := range dict {
		_words = append(_words, r)
	}
	sort.Sort(ByRune(_words))
	nextCode := 1
	for _, r := range _words {
		dict[r] = nextCode
		nextCode++
	}
	return &WordCodeDict{nextCode: nextCode, dict: dict}
}

func (d *WordCodeDict) Code(word rune) int {
	if code, ok := d.dict[word]; ok {
		return code
	}
	// 返回一个非法值
	return d.nextCode
}

type AhoCorasickDoubleArrayTrie struct {
	base   []int
	check  []int
	used   []bool
	size   int
	allocSize int
	nextCheckPos int
	keySize int
	progress int
	rootState *State
	// length of every key
	l       []int
	value  []interface{}
	// fail table of the Aho Corasick automata
	fail   []int
	// output table of the Aho Corasick automata
	output  [][]int
	wordCodeDict *WordCodeDict
}

func NewAhoCorasickDoubleArrayTrie() *AhoCorasickDoubleArrayTrie {
	return &AhoCorasickDoubleArrayTrie{rootState: NewState(0, 0)}
}

func (act *AhoCorasickDoubleArrayTrie) ParseText(key string) []*Hit {
	var position int = 1
	var curState int = 0
	var collectedEmits = NewListHit()
	text := NewWord(key)
	for i := 0; i < text.Size(); i++ {
		c := act.wordCodeDict.Code(text.GetRune(i))
		curState = act.getState(curState, c)
		act.storeEmits(position, curState, collectedEmits)
		position++
	}
	return collectedEmits.ListArray()
}

func (act *AhoCorasickDoubleArrayTrie) ParseTextWithIter(key string, iter IHit) {
	var position int = 1
	var curState int = 0
	text := NewWord(key)
	for i := 0; i < text.Size(); i++ {
		c := act.wordCodeDict.Code(text.GetRune(i))
		curState = act.getState(curState, c)
		hitArray := act.output[curState]
		for _, hit := range hitArray {
			iter(position - act.l[hit], position, act.value[hit])
		}
		position++
	}
}

func (act *AhoCorasickDoubleArrayTrie) Matches(key string) bool {
	var curState int = 0
	text := NewWord(key)
	for i := 0; i < text.Size(); i++ {
		c := act.wordCodeDict.Code(text.GetRune(i))
		curState = act.getState(curState, c)
		hitArray := act.output[curState]
		if hitArray != nil {
			return true
		}
	}
	return false
}

func (act *AhoCorasickDoubleArrayTrie) FindFirst(key string) *Hit {
	var position int = 1
	var curState int = 0
	text := NewWord(key)
	for i := 0; i < text.Size(); i++ {
		c := act.wordCodeDict.Code(text.GetRune(i))
		curState = act.getState(curState, c)
		hitArray := act.output[curState]
		for _, hit := range hitArray {
			return NewHit(position - act.l[hit], position, act.value[hit])
		}
		position++
	}
	return nil
}

func (act *AhoCorasickDoubleArrayTrie) Get(key string) interface{} {
	index := act.ExactMatchSearch(key)
	if index >= 0 {
		return act.value[index]
	}
	return nil
}

func (act *AhoCorasickDoubleArrayTrie) Build(kvs *StringTreeMap) {
	act.value = kvs.Values()
	act.l = make([]int, len(act.value))

	// 构建二分trie树
	act.addAllKeyWord(kvs.Keys())
	act.buildDoubleArrayTrie(len(act.value))
	act.used = nil
	act.constructFailureStates()
	act.rootState.String()
	act.rootState = nil
	act.loseWeight()
}

func (act *AhoCorasickDoubleArrayTrie) Dump() {
	fmt.Println("base: ", act.base)
	fmt.Println("check: ", act.check)
	fmt.Println("fail: ", act.fail)
	fmt.Println("output: ",act.output)
}

func(act *AhoCorasickDoubleArrayTrie) getState(currentState int, code int) int {
	newCurrentState := act.transitionWithRoot(currentState, code)  // 先按success跳转
	for newCurrentState == -1 {
		currentState = act.fail[currentState]
		newCurrentState = act.transitionWithRoot(currentState, code)
	}
	return newCurrentState
}

func(act *AhoCorasickDoubleArrayTrie) storeEmits(position int, currectState int, collectEmits *ListHit) {
	hitArray := act.output[currectState]
	for _, hit := range hitArray {
		collectEmits.add(NewHit(position-act.l[hit], position, act.value[hit]))
	}
}

func(act *AhoCorasickDoubleArrayTrie)transition(current int, c int) int {
	b := current
	p := b+c+1
	if b == act.check[p] {
		b = act.base[p]
	} else {
		return -1
	}
	p = b
	return p
}

func (act *AhoCorasickDoubleArrayTrie) transitionWithRoot(nodePos int, c int) int {
	b := act.base[nodePos]
	p := b + c + 1
	if b != act.check[p] {
		if nodePos == 0 {
			return 0
		}
		return -1
	}
	return p
}

func (act *AhoCorasickDoubleArrayTrie)resize(newSize int) int {
	base2 := make([]int, newSize)
	check2 := make([]int, newSize)
	used2 := make([]bool, newSize)

	if act.allocSize > 0 {
		copy(base2, act.base)
		copy(check2, act.check)
		copy(used2, act.used)
	}
	act.base = base2
	act.check = check2
	act.used = used2
	act.allocSize = newSize
	return newSize
}

func (act *AhoCorasickDoubleArrayTrie) fetch(parent *State, siblings *ListState) int {
	if parent.isAcceptable() {
		s := NewState((parent.getDepth() + 1) * -1, 0)
		s.addEmit(parent.getLargestValueId())
		siblings.add(s)
	}
	for _, v := range parent.getSuccess() {
		siblings.add(v)
	}
	return siblings.size()
}

func (act *AhoCorasickDoubleArrayTrie) addKeyWord(text *Word, index int) {
	curState := act.rootState
	for _, c := range text.GetRunes() {
		curState = curState.addState(c, act.wordCodeDict.Code(c) + 1)
	}
	curState.addEmit(index)
	act.l[index] = text.Size()
}

func (act *AhoCorasickDoubleArrayTrie) addAllKeyWord(keyWordSet []string) {
	var words []*Word
	for _, keyword := range keyWordSet {
		text := NewWord(keyword)
		words = append(words, text)
	}
	act.wordCodeDict = NewWordCodeDict(words)
	for i, keyword := range words {
		words = append(words, keyword)
		act.addKeyWord(keyword, i)
	}
}

func (act *AhoCorasickDoubleArrayTrie) constructFailureStates() {
	act.fail = make([]int, act.size + 1)
	act.fail[1] = act.base[0]
	act.output = make([][]int, act.size+1)
	
	queue := NewArrayStateQueue()
	for _, depthOneState := range act.rootState.getStates() {
		depthOneState.setFailure(act.rootState, act.fail)
		queue.add(depthOneState)
		act.constructOutput(depthOneState)
	}
	for !queue.isEmpty() {
		curState := queue.remove()
		for _, c := range curState.getTransitions() {
			targetState := curState.nextState(c)
			queue.add(targetState)
			
			traceFailureState := curState.failure()
			for traceFailureState.nextState(c) == nil {
				traceFailureState = traceFailureState.failure()
			}
			
			newFailureState := traceFailureState.nextState(c)
			targetState.setFailure(newFailureState, act.fail)
			targetState.addEmits(newFailureState.emit())
			act.constructOutput(targetState)
		}
	}
}

func (act *AhoCorasickDoubleArrayTrie) constructOutput(targetState *State) {
	emit := targetState.emit()
	if emit == nil || len(emit) == 0 {
		return
	}
	output := make([]int, len(emit))
	for i := 0; i < len(output); i++ {
		output[i] = emit[i]
	}
	act.output[targetState.getIndex()] = output
}

func (act *AhoCorasickDoubleArrayTrie) insert(siblings *ListState) int {
	begin := 0
	nonzero_num := 0
	first := 0
	var pos int
	if siblings.get(0).code + 1 > act.nextCheckPos {
		pos = siblings.get(0).code + 1
	} else {
		pos = act.nextCheckPos
	}
	pos -= 1

	if act.allocSize <= pos {
		act.resize(pos + 1)
	}
	OUTER:
	// 此循环体的目标是找出满足base[begin + a1...an]==0, check[begin + a1...an]==0的n个空闲空间,a1...an是siblings中的n个节点
	for {
		pos++

		if act.allocSize <= pos {
			act.resize(pos+1)
		}
		if act.check[pos] != 0 {
			nonzero_num++
			continue
		} else if first == 0 {
			act.nextCheckPos = pos
			first = 1
		}
		begin = pos - siblings.get(0).code
		if act.allocSize <= (begin + siblings.get(siblings.size() - 1).code) {
			// progress can be zero
			var l float64
			tmp_l := 1.0 * float64(act.keySize) / float64(act.progress + 1)
			if 1.05 > tmp_l {
				l = 1.05
			} else {
				l = tmp_l
			}
			act.resize(int(float64(act.allocSize) * l))
		}
        // 这个位置已经被使用了
		if act.used[begin] {
			continue
		}

		// 检查是否存在冲突
		for i := 0; i < siblings.size(); i++ {
			if act.base[begin + siblings.get(i).code] != 0 {
				continue OUTER
			}
			if act.check[begin + siblings.get(i).code] != 0 {
				continue OUTER
			}
		}
		// 找到一个没有冲突的位置
		break
	}
	if 1.0 * float64(nonzero_num) / float64(pos - act.nextCheckPos + 1) >= 0.95 {
		act.nextCheckPos = pos
	}
    // 标记位置被占用
	act.used[begin] = true
	tmp_size := begin + siblings.get(siblings.size() - 1).code + 1
	// 更新 tire的size
	if act.size < tmp_size {
		act.size = tmp_size
	}

	// base[s] + c = t
	// check[t] = s
	for i := 0; i < siblings.size(); i++ {
		act.check[begin + siblings.get(i).code] = begin
	}

	// 计算所有子节点的base
	for i := 0; i < siblings.size(); i++ {
		new_siblings := NewListState()
		//// 一个词的终止且不为其他词的前缀，其实就是叶子节点
		if act.fetch(siblings.get(i), new_siblings) == 0 {
			act.base[begin+siblings.get(i).code] = siblings.get(i).getLargestValueId() * (-1) - 1
			act.progress++
		} else {
			h := act.insert(new_siblings)
			act.base[begin+siblings.get(i).code] = h
		}
		siblings.get(i).setIndex(begin + siblings.get(i).code)
	}

	return begin
}

func (act *AhoCorasickDoubleArrayTrie) loseWeight() {
	base2 := make([]int, act.size)
	check2 := make([]int, act.size)

	if act.allocSize > 0 {
		copy(base2, act.base)
		copy(check2, act.check)
	}
	act.base = base2
	act.check = check2
	act.allocSize = act.size
	return
}

func (act *AhoCorasickDoubleArrayTrie) GetSize() int {
	return act.size
}

func (act *AhoCorasickDoubleArrayTrie) GetNonzeroSize() int {
	result := 0
	for i := 0; i< act.size; i++ {
		if act.check[i] != 0 {
			result++
		}
	}
	return result
}

func (act *AhoCorasickDoubleArrayTrie) buildDoubleArrayTrie(keySize int) {
	act.keySize = keySize
	act.progress = 0

	// 32个双字节
	act.resize(65536 * 32)

	act.base[0] = 1
	act.nextCheckPos = 0

	root_state := act.rootState

	siblings := NewListState()
	act.fetch(root_state, siblings)
	act.insert(siblings)
	return
}

func (act *AhoCorasickDoubleArrayTrie) ExactMatchSearch(key string) int {
	return act.ExactMatchSearchAdvanced(key, 0, 0, 0)
}

func (act *AhoCorasickDoubleArrayTrie) ExactMatchSearchAdvanced(key string, pos int, length int, nodePos int) int {
	word := NewWord(key)
	if length <= 0 {
		length = word.Size()
	}
	if nodePos <= 0 {
		nodePos = 0
	}

	var result = -1

	keyChars := word.GetRunes()
	b := act.base[nodePos]
	var p int
	for i := pos; i < length; i++ {
		p = b + act.wordCodeDict.Code(keyChars[i]) + 1
		if b == act.check[p] {
			b = act.base[p]
		} else {
			return result
		}
	}

	p = b
	n := act.base[p]

	if b == act.check[p] && n < 0 {
		result = n * (-1) - 1
	}
	return result
}


type ByRune []rune
func (a ByRune) Len() int           { return len(a) }
func (a ByRune) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByRune) Less(i, j int) bool { return a[i] < a[j] }


