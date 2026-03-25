package tui

import (
	"fmt"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// treeNode represents a file or folder in the repo tree.
type treeNode struct {
	name     string
	path     string // full relative path
	isDir    bool
	children []*treeNode
	selected bool
	expanded bool
	depth    int
}

// PickerModel is the state for the file/folder picker.
type PickerModel struct {
	root       *treeNode
	flatList   []*treeNode // visible items in current expanded state
	cursor     int
	folderMode bool
	scrollOff  int
	viewHeight int
}

// NewPickerModel builds a tree from a flat list of file paths.
func NewPickerModel(files []string, folderMode bool) PickerModel {
	root := &treeNode{name: "/", path: "", isDir: true, expanded: true, depth: -1}
	nodeMap := map[string]*treeNode{"": root}

	// Build tree
	for _, f := range files {
		parts := strings.Split(f, "/")
		for i := range parts {
			partialPath := strings.Join(parts[:i+1], "/")
			if _, exists := nodeMap[partialPath]; exists {
				continue
			}

			isDir := i < len(parts)-1
			node := &treeNode{
				name:     parts[i],
				path:     partialPath,
				isDir:    isDir,
				expanded: false,
				depth:    i,
			}

			parentPath := ""
			if i > 0 {
				parentPath = strings.Join(parts[:i], "/")
			}
			parent := nodeMap[parentPath]
			parent.children = append(parent.children, node)
			nodeMap[partialPath] = node
		}
	}

	// Sort children alphabetically, directories first
	sortTree(root)

	// If folder mode, expand first level
	if folderMode {
		root.expanded = true
	} else {
		root.expanded = true
	}

	pm := PickerModel{
		root:       root,
		folderMode: folderMode,
		viewHeight: 20,
	}
	pm.rebuildFlatList()
	return pm
}

func sortTree(node *treeNode) {
	sort.Slice(node.children, func(i, j int) bool {
		// Dirs first
		if node.children[i].isDir != node.children[j].isDir {
			return node.children[i].isDir
		}
		return node.children[i].name < node.children[j].name
	})
	for _, child := range node.children {
		if child.isDir {
			sortTree(child)
		}
	}
}

// rebuildFlatList creates the visible item list based on expansion state.
func (pm *PickerModel) rebuildFlatList() {
	pm.flatList = nil
	pm.buildFlat(pm.root)
}

func (pm *PickerModel) buildFlat(node *treeNode) {
	for _, child := range node.children {
		if pm.folderMode && !child.isDir {
			continue
		}
		pm.flatList = append(pm.flatList, child)
		if child.isDir && child.expanded {
			pm.buildFlat(child)
		}
	}
}

// getSelectedFiles returns all selected file paths.
func (pm *PickerModel) getSelectedFiles(allFiles []string) []string {
	if pm.folderMode {
		return pm.getSelectedFolderFiles(allFiles)
	}
	var result []string
	pm.collectSelected(pm.root, &result)
	return result
}

func (pm *PickerModel) collectSelected(node *treeNode, result *[]string) {
	for _, child := range node.children {
		if child.isDir {
			pm.collectSelected(child, result)
		} else if child.selected {
			*result = append(*result, child.path)
		}
	}
}

func (pm *PickerModel) getSelectedFolderFiles(allFiles []string) []string {
	var selectedDirs []string
	pm.collectSelectedDirs(pm.root, &selectedDirs)

	var result []string
	for _, f := range allFiles {
		for _, dir := range selectedDirs {
			if strings.HasPrefix(f, dir+"/") || f == dir {
				result = append(result, f)
				break
			}
		}
	}
	return result
}

func (pm *PickerModel) collectSelectedDirs(node *treeNode, result *[]string) {
	for _, child := range node.children {
		if child.isDir {
			if child.selected {
				*result = append(*result, child.path)
			}
			pm.collectSelectedDirs(child, result)
		}
	}
}

// toggle selection of the item at cursor.
func (pm *PickerModel) toggleCurrent() {
	if pm.cursor < 0 || pm.cursor >= len(pm.flatList) {
		return
	}
	node := pm.flatList[pm.cursor]

	if pm.folderMode {
		if node.isDir {
			node.selected = !node.selected
		}
	} else {
		if node.isDir {
			// Toggle all children recursively
			newState := !pm.allChildrenSelected(node)
			pm.setChildrenSelected(node, newState)
			node.selected = newState
		} else {
			node.selected = !node.selected
		}
	}
}

func (pm *PickerModel) allChildrenSelected(node *treeNode) bool {
	for _, child := range node.children {
		if child.isDir {
			if !pm.allChildrenSelected(child) {
				return false
			}
		} else if !child.selected {
			return false
		}
	}
	return true
}

func (pm *PickerModel) setChildrenSelected(node *treeNode, state bool) {
	for _, child := range node.children {
		child.selected = state
		if child.isDir {
			pm.setChildrenSelected(child, state)
		}
	}
}

// toggleExpand expands or collapses a directory.
func (pm *PickerModel) toggleExpand() {
	if pm.cursor < 0 || pm.cursor >= len(pm.flatList) {
		return
	}
	node := pm.flatList[pm.cursor]
	if node.isDir {
		node.expanded = !node.expanded
		pm.rebuildFlatList()
		// Ensure cursor is still valid
		if pm.cursor >= len(pm.flatList) {
			pm.cursor = len(pm.flatList) - 1
		}
	}
}

// selectAll toggles selection of all visible items.
func (pm *PickerModel) selectAll() {
	allSelected := true
	for _, node := range pm.flatList {
		if pm.folderMode && !node.isDir {
			continue
		}
		if !node.selected {
			allSelected = false
			break
		}
	}

	newState := !allSelected
	for _, node := range pm.flatList {
		if pm.folderMode && !node.isDir {
			continue
		}
		node.selected = newState
		if node.isDir && !pm.folderMode {
			pm.setChildrenSelected(node, newState)
		}
	}
}

// countSelected returns how many items are selected.
func (pm *PickerModel) countSelected() int {
	count := 0
	for _, node := range pm.flatList {
		if node.selected {
			count++
		}
	}
	return count
}

// ---------- TUI integration ----------

func (m Model) updatePicker(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.picker.viewHeight = msg.Height - 10
		if m.picker.viewHeight < 5 {
			m.picker.viewHeight = 5
		}

	case tea.KeyMsg:
		switch {
		case msg.String() == "up" || msg.String() == "k":
			if m.picker.cursor > 0 {
				m.picker.cursor--
				m.picker.ensureVisible()
			}

		case msg.String() == "down" || msg.String() == "j":
			if m.picker.cursor < len(m.picker.flatList)-1 {
				m.picker.cursor++
				m.picker.ensureVisible()
			}

		case msg.String() == " ":
			m.picker.toggleCurrent()

		case msg.String() == "tab":
			m.picker.toggleExpand()

		case msg.String() == "a":
			m.picker.selectAll()

		case msg.String() == "enter":
			selected := m.picker.getSelectedFiles(m.files)
			if len(selected) == 0 {
				return m, nil
			}
			m.selectedFiles = selected
			m.screen = ScreenFetching
			return m, tea.Batch(m.spinner.Tick, fetchCmd(m.link.Path, selected))

		case msg.String() == "esc":
			m.screen = ScreenModeSelect
			return m, nil

		case msg.String() == "q":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (pm *PickerModel) ensureVisible() {
	if pm.cursor < pm.scrollOff {
		pm.scrollOff = pm.cursor
	}
	if pm.cursor >= pm.scrollOff+pm.viewHeight {
		pm.scrollOff = pm.cursor - pm.viewHeight + 1
	}
}

func (m Model) viewPicker() string {
	var b strings.Builder
	pm := &m.picker

	modeLabel := "Select Files"
	if pm.folderMode {
		modeLabel = "Select Folders"
	}
	b.WriteString(subtitleStyle.Render(modeLabel))
	b.WriteString("\n")

	selectedCount := pm.countSelected()
	b.WriteString(dimStyle.Render(fmt.Sprintf(
		"%s selected  •  %s total",
		fmt.Sprintf("%d", selectedCount),
		pluralize(len(pm.flatList), "item", "items"),
	)))
	b.WriteString("\n\n")

	if len(pm.flatList) == 0 {
		b.WriteString(dimStyle.Render("No items to display."))
		return b.String()
	}

	// Render visible window
	end := pm.scrollOff + pm.viewHeight
	if end > len(pm.flatList) {
		end = len(pm.flatList)
	}

	for i := pm.scrollOff; i < end; i++ {
		node := pm.flatList[i]
		isCurrent := i == pm.cursor

		// Indentation
		indent := strings.Repeat("  ", node.depth)

		// Selection indicator
		var check string
		if node.selected {
			check = checkedStyle.Render("[✓]")
		} else {
			check = uncheckedStyle.Render("[ ]")
		}

		// Dir expand indicator
		var icon string
		if node.isDir {
			if node.expanded {
				icon = "▾ 📁 "
			} else {
				icon = "▸ 📁 "
			}
		} else {
			icon = "  📄 "
		}

		// Name
		name := node.name
		if isCurrent {
			name = selectedStyle.Render(name)
		} else {
			name = normalStyle.Render(name)
		}

		cursor := "  "
		if isCurrent {
			cursor = selectedStyle.Render("▸ ")
		}

		line := fmt.Sprintf("%s%s%s %s%s", cursor, indent, check, icon, name)
		b.WriteString(line)
		b.WriteString("\n")
	}

	// Scroll indicator
	if len(pm.flatList) > pm.viewHeight {
		b.WriteString("\n")
		progress := float64(pm.scrollOff) / float64(len(pm.flatList)-pm.viewHeight)
		bar := renderScrollBar(progress, pm.viewHeight, len(pm.flatList))
		b.WriteString(dimStyle.Render(bar))
	}

	return b.String()
}

func renderScrollBar(progress float64, viewHeight, totalItems int) string {
	pos := int(progress * 10)
	if pos > 9 {
		pos = 9
	}
	bar := "["
	for i := 0; i < 10; i++ {
		if i == pos {
			bar += "█"
		} else {
			bar += "░"
		}
	}
	bar += fmt.Sprintf("] %d/%d", clamp(int(progress*float64(totalItems))+1, 1, totalItems), totalItems)
	return bar
}
