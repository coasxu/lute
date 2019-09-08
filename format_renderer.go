// Lute - A structured markdown engine.
// Copyright (c) 2019-present, b3log.org
//
// Lute is licensed under the Mulan PSL v1.
// You can use this software according to the terms and conditions of the Mulan PSL v1.
// You may obtain a copy of Mulan PSL v1 at:
//     http://license.coscl.org.cn/MulanPSL
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR
// PURPOSE.
// See the Mulan PSL v1 for more details.

package lute

import (
	"bytes"
	"strconv"
)

// newFormatRenderer 创建一个格式化渲染器。
func (lute *Lute) newFormatRenderer(treeRoot *Node) (ret *Renderer) {
	ret = &Renderer{rendererFuncs: map[int]RendererFunc{}, option: lute.options, treeRoot: treeRoot}

	// 注册 CommonMark 渲染函数

	ret.rendererFuncs[NodeDocument] = ret.renderDocumentMarkdown
	ret.rendererFuncs[NodeParagraph] = ret.renderParagraphMarkdown
	ret.rendererFuncs[NodeText] = ret.renderTextMarkdown
	ret.rendererFuncs[NodeCodeSpan] = ret.renderCodeSpanMarkdown
	ret.rendererFuncs[NodeCodeBlock] = ret.renderCodeBlockMarkdown
	ret.rendererFuncs[NodeEmphasis] = ret.renderEmphasisMarkdown
	ret.rendererFuncs[NodeStrong] = ret.renderStrongMarkdown
	ret.rendererFuncs[NodeBlockquote] = ret.renderBlockquoteMarkdown
	ret.rendererFuncs[NodeHeading] = ret.renderHeadingMarkdown
	ret.rendererFuncs[NodeList] = ret.renderListMarkdown
	ret.rendererFuncs[NodeListItem] = ret.renderListItemMarkdown
	ret.rendererFuncs[NodeThematicBreak] = ret.renderThematicBreakMarkdown
	ret.rendererFuncs[NodeHardBreak] = ret.renderHardBreakMarkdown
	ret.rendererFuncs[NodeSoftBreak] = ret.renderSoftBreakMarkdown
	ret.rendererFuncs[NodeHTMLBlock] = ret.renderHTMLMarkdown
	ret.rendererFuncs[NodeInlineHTML] = ret.renderInlineHTMLMarkdown
	ret.rendererFuncs[NodeLink] = ret.renderLinkMarkdown
	ret.rendererFuncs[NodeImage] = ret.renderImageMarkdown

	// 注册 GFM 渲染函数

	ret.rendererFuncs[NodeStrikethrough] = ret.renderStrikethroughMarkdown
	ret.rendererFuncs[NodeTaskListItemMarker] = ret.renderTaskListItemMarkerMarkdown
	ret.rendererFuncs[NodeTable] = ret.renderTableMarkdown
	ret.rendererFuncs[NodeTableHead] = ret.renderTableHeadMarkdown
	ret.rendererFuncs[NodeTableRow] = ret.renderTableRowMarkdown
	ret.rendererFuncs[NodeTableCell] = ret.renderTableCellMarkdown

	// Emoji 渲染函数

	ret.rendererFuncs[NodeEmojiUnicode] = ret.renderEmojiUnicodeMarkdown
	ret.rendererFuncs[NodeEmojiImg] = ret.renderEmojiImgMarkdown

	return
}

func (r *Renderer) renderEmojiImgMarkdown(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.write(node.emojiImgAlias)
	}
	return WalkContinue, nil
}

func (r *Renderer) renderEmojiUnicodeMarkdown(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.write(node.tokens)
	}
	return WalkContinue, nil
}

// TODO: 表的格式化应该按最宽的单元格对齐内容

func (r *Renderer) renderTableCellMarkdown(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.writeByte(itemPipe)
	}
	return WalkContinue, nil
}

func (r *Renderer) renderTableRowMarkdown(node *Node, entering bool) (WalkStatus, error) {
	if !entering {
		r.writeString("|\n")
	}
	return WalkContinue, nil
}

func (r *Renderer) renderTableHeadMarkdown(node *Node, entering bool) (WalkStatus, error) {
	if !entering {
		r.writeString("|\n")
		table := node.parent
		for i := 0; i < len(table.tableAligns); i++ {
			align := table.tableAligns[i]
			switch align {
			case 0:
				r.writeString("|---")
			case 1:
				r.writeString("|:---")
			case 2:
				r.writeString("|:---:")
			case 3:
				r.writeString("|---:")
			}
		}
		r.writeString("|\n")
	}
	return WalkContinue, nil
}

func (r *Renderer) renderTableMarkdown(node *Node, entering bool) (WalkStatus, error) {
	if !entering {
		r.newline()
		if !r.isLastNode(r.treeRoot, node) {
			r.writeByte(itemNewline)
		}
	}
	return WalkContinue, nil
}

func (r *Renderer) renderStrikethroughMarkdown(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.write(bytes.Repeat(items{node.strongEmDelMarker}, node.strongEmDelMarkenLen))
	} else {
		r.write(bytes.Repeat(items{node.strongEmDelMarker}, node.strongEmDelMarkenLen))
	}
	return WalkContinue, nil
}

func (r *Renderer) renderImageMarkdown(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.writeString("![")
		r.write(node.firstChild.tokens)
		r.writeString("](")
		r.write(node.destination)
		if nil != node.title {
			r.writeString(" \"")
			r.write(node.title)
			r.writeByte(itemDoublequote)
		}
		r.writeByte(itemCloseParen)
	}
	return WalkContinue, nil
}

func (r *Renderer) renderLinkMarkdown(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.writeString("[")
		if nil != node.firstChild {
			// FIXME: 未解决链接嵌套，另外还需要考虑链接引用定义
			r.write(node.firstChild.tokens)
		}
		r.writeString("](")
		r.write(node.destination)
		if nil != node.title {
			r.writeString(" \"")
			r.write(node.title)
			r.writeByte(itemDoublequote)
		}
		r.writeByte(itemCloseParen)
	}
	return WalkContinue, nil
}

func (r *Renderer) renderHTMLMarkdown(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.newline()
		r.write(node.tokens)
		r.newline()
	}
	return WalkContinue, nil
}

func (r *Renderer) renderInlineHTMLMarkdown(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.write(node.tokens)
	}
	return WalkContinue, nil
}

func (r *Renderer) renderDocumentMarkdown(node *Node, entering bool) (WalkStatus, error) {
	if !entering {
		r.writeByte(itemNewline)
	}
	return WalkContinue, nil
}

func (r *Renderer) renderParagraphMarkdown(node *Node, entering bool) (WalkStatus, error) {
	listPadding := 0
	inTightList := false
	lastListItemLastPara := false
	if parent := node.parent; nil != parent {
		if NodeListItem == parent.typ { // ListItem.Paragraph
			listItem := parent

			// 必须通过列表（而非列表项）上的紧凑标识判断，因为在设置该标识时仅设置了 List.tight
			// 设置紧凑标识的具体实现可参考函数 List.Finalize()
			inTightList = listItem.parent.tight

			firstPara := listItem.firstChild
			if 3 != listItem.listData.typ { // 普通列表
				if firstPara != node {
					listPadding = listItem.padding
				}
			} else { // 任务列表
				if firstPara.next != node { // 任务列表要跳过 TaskListItemMarker 即 [X]
					listPadding = listItem.padding
				}
			}

			nextItem := listItem.next
			if nil == nextItem {
				nextPara := node.next
				lastListItemLastPara = nil == nextPara
			}
		}
	}

	if entering {
		r.write(bytes.Repeat(items{itemSpace}, listPadding))
	} else {
		r.newline()
		isLastNode := r.isLastNode(r.treeRoot, node)
		if !isLastNode && (!inTightList || (lastListItemLastPara)) {
			r.writeByte(itemNewline)
		}
	}
	return WalkContinue, nil
}

func (r *Renderer) renderTextMarkdown(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		if typ := node.parent.typ; NodeLink != typ && NodeImage != typ {
			r.write(escapeHTML(node.tokens))
		}
	}
	return WalkContinue, nil
}

func (r *Renderer) renderCodeSpanMarkdown(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.writeByte(itemBacktick)
		if 1 < node.codeMarkerLen {
			r.writeByte(itemBacktick)
			r.writeByte(itemSpace)
		}
		r.write(node.tokens)
		return WalkSkipChildren, nil
	}

	if 1 < node.codeMarkerLen {
		r.writeByte(itemSpace)
		r.writeByte(itemBacktick)
	}
	r.writeByte(itemBacktick)
	return WalkContinue, nil
}

func (r *Renderer) renderCodeBlockMarkdown(node *Node, entering bool) (WalkStatus, error) {
	if !node.isFencedCodeBlock {
		node.codeBlockFenceLen = 3
	}
	if entering {
		listPadding := 0
		if grandparent := node.parent.parent; nil != grandparent {
			if NodeList == grandparent.typ { // List.ListItem.CodeBlock
				if node.parent.firstChild != node {
					listPadding = grandparent.padding
				}
			}
		}

		r.newline()
		if 0 < listPadding {
			r.write(bytes.Repeat(items{itemSpace}, listPadding))
		}
		r.write(bytes.Repeat(items{itemBacktick}, node.codeBlockFenceLen))
		r.write(node.codeBlockInfo)
		r.writeByte(itemNewline)
		if 0 < listPadding {
			lines := bytes.Split(node.tokens, items{itemNewline})
			length := len(lines)
			for i, line := range lines {
				r.write(bytes.Repeat(items{itemSpace}, listPadding))
				r.write(line)
				if i < length-1 {
					r.writeByte(itemNewline)
				}
			}
		} else {
			r.write(node.tokens)
		}
		return WalkSkipChildren, nil
	}

	r.write(bytes.Repeat(items{itemBacktick}, node.codeBlockFenceLen))
	r.newline()
	if !r.isLastNode(r.treeRoot, node) {
		r.writeByte(itemNewline)
	}
	return WalkContinue, nil
}

func (r *Renderer) renderEmphasisMarkdown(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.write(bytes.Repeat(items{node.strongEmDelMarker}, node.strongEmDelMarkenLen))
	} else {
		r.write(bytes.Repeat(items{node.strongEmDelMarker}, node.strongEmDelMarkenLen))
	}
	return WalkContinue, nil
}

func (r *Renderer) renderStrongMarkdown(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.write(bytes.Repeat(items{node.strongEmDelMarker}, node.strongEmDelMarkenLen))
	} else {
		r.write(bytes.Repeat(items{node.strongEmDelMarker}, node.strongEmDelMarkenLen))
	}
	return WalkContinue, nil
}

func (r *Renderer) renderBlockquoteMarkdown(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.newline()
		r.writeString("> ") // 带个空格更好一些
	} else {
		r.newline()
	}
	return WalkContinue, nil
}

func (r *Renderer) renderHeadingMarkdown(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.write(bytes.Repeat(items{itemCrosshatch}, node.headingLevel)) // 统一使用 ATX 标题，不使用 Setext 标题
		r.writeByte(itemSpace)
	} else {
		r.newline()
		r.writeByte(itemNewline)
	}
	return WalkContinue, nil
}

func (r *Renderer) renderListMarkdown(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.newline()
		r.listDepth++
		if 1 < r.listDepth {
			lastList := r.listStack[len(r.listStack)-1] // 栈顶是上一个列表节点
			r.listIndent += len(lastList.marker) + 1
			if 1 == lastList.listData.typ {
				r.listIndent++
			}
		}
		r.listStack = append(r.listStack, node) // 入栈
	} else {
		r.newline()
		r.listStack = r.listStack[:len(r.listStack)-1] // 出栈
		if 0 < len(r.listStack) {
			lastList := r.listStack[len(r.listStack)-1]
			r.listIndent -= + len(lastList.marker) + 1
			if 1 == lastList.listData.typ {
				r.listIndent--
			}
		}
		r.listDepth--
	}
	return WalkContinue, nil
}

func (r *Renderer) renderListItemMarkdown(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.newline()
		if 1 < r.listDepth {
			r.write(bytes.Repeat(items{itemSpace}, r.listIndent))
		}
		if 1 == node.listData.typ {
			r.writeString(strconv.Itoa(node.num) + ".")
		} else {
			r.write(node.marker)
		}
		r.writeByte(itemSpace)
	}
	return WalkContinue, nil
}

func (r *Renderer) renderTaskListItemMarkerMarkdown(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.writeByte(itemOpenBracket)
		if node.taskListItemChecked {
			r.writeByte('X')
		} else {
			r.writeByte(itemSpace)
		}
		r.writeByte(itemCloseBracket)
	}
	return WalkContinue, nil
}

func (r *Renderer) renderThematicBreakMarkdown(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.newline()
		r.writeString("---")
		r.newline()
	}
	return WalkSkipChildren, nil
}

func (r *Renderer) renderHardBreakMarkdown(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		if !r.option.SoftBreak2HardBreak {
			r.writeString("\\\n")
		} else {
			r.writeByte(itemNewline)
		}
	}
	return WalkSkipChildren, nil
}

func (r *Renderer) renderSoftBreakMarkdown(node *Node, entering bool) (WalkStatus, error) {
	if entering {
		r.newline()
	}
	return WalkSkipChildren, nil
}

func (r *Renderer) isLastNode(treeRoot, node *Node) bool {
	if treeRoot == node {
		return true
	}

	if NodeDocument == node.parent.typ {
		return treeRoot.lastChild == node
	}

	var n *Node
	for n = node.parent; ; n = n.parent {
		if NodeDocument == n.parent.typ {
			break
		}
	}
	return treeRoot.lastChild == n
}
