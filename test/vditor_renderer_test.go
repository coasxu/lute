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

// +build javascript

package test

import (
	"testing"

	"github.com/b3log/lute"
)

var vditorRendererTests = []parseTest{
	//{"20", "> *foo\n> bar*\n", "<blockquote data-ntype=\"4\" data-mtype=\"1\"><span class=\"node\"><span class=\"marker\" data-ntype=\"5\" data-mtype=\"2\" data-caret=\"start\" data-caretoffset=\"0\">&gt;</span></span><p data-ntype=\"1\" data-mtype=\"0\"><span class=\"node\" data-ntype=\"12\" data-mtype=\"2\"><span class=\"marker\" data-ntype=\"13\" data-mtype=\"2\">*</span><em data-ntype=\"12\" data-mtype=\"2\"><span data-ntype=\"11\" data-mtype=\"2\">foo</span><span data-ntype=\"24\" data-mtype=\"2\"><br><span class=\"newline\" data-ntype=\"24\" data-mtype=\"2\" />\n</span><span data-ntype=\"11\" data-mtype=\"2\">bar</span></em><span class=\"marker\" data-ntype=\"14\" data-mtype=\"2\">*</span></span></p><span><br><span class=\"newline\">\n</span><span class=\"newline\">\n</span></span></blockquote>"},
	//{"14", "1\n\n2\n", "<p data-ntype=\"1\" data-mtype=\"0\"><span data-ntype=\"1 data-mtype=\"2\" data-cso=\"0\" data-ceo=\"0\"></span><span class=\"newline\">\n</span><span class=\"newline\">\n</span></p>"},
	{"13", "", "<p data-ntype=\"1\" data-mtype=\"0\"><span data-ntype=\"1 data-mtype=\"2\" data-cso=\"0\" data-ceo=\"0\"></span><span class=\"newline\">\n</span><span class=\"newline\">\n</span></p>"},
	{"12", "foo\nbar\n", "<p data-ntype=\"1\" data-mtype=\"0\">foo\nbar<span class=\"newline\">\n\n</span></p>"},
	{"11", "> # foo\n", "<div class=\"node node--block node--expand\" data-ntype=\"5\" data-mtype=\"1\"><span class=\"marker\" data-cso=\"2\" data-ceo=\"2\">&gt; </span><blockquote data-ntype=\"5\" data-mtype=\"1\"><h1 class=\"node\" data-ntype=\"2\" data-mtype=\"0\"><span class=\"marker\"># </span>foo</h1></blockquote></div>"},
	{"10", "> #\n", "<div class=\"node node--block node--expand\" data-ntype=\"5\" data-mtype=\"1\"><span class=\"marker\" data-cso=\"2\" data-ceo=\"2\">&gt; </span><blockquote data-ntype=\"5\" data-mtype=\"1\"><h1 class=\"node\" data-ntype=\"2\" data-mtype=\"0\"><span class=\"marker\">#\n</span></h1></blockquote></div>"},
	{"9", "> ## foo\n", "<div class=\"node node--block node--expand\" data-ntype=\"5\" data-mtype=\"1\"><span class=\"marker\" data-cso=\"2\" data-ceo=\"2\">&gt; </span><blockquote data-ntype=\"5\" data-mtype=\"1\"><h2 class=\"node\" data-ntype=\"2\" data-mtype=\"0\"><span class=\"marker\">## </span>foo</h2></blockquote></div>"},
	{"8", "-- ---\n", "<hr data-ntype=\"4\" data-mtype=\"0\" data-cso=\"2\" data-ceo=\"2\" />"},
	{"7", "__foo__\n", "<p data-ntype=\"1\" data-mtype=\"0\"><span class=\"node node--expand\" data-ntype=\"18\" data-mtype=\"2\"><span class=\"marker\" data-cso=\"2\" data-ceo=\"2\">__</span><strong data-ntype=\"18\" data-mtype=\"2\">foo</strong><span class=\"marker\">__</span></span><span class=\"newline\">\n\n</span></p>"},
	{"6", "`foo`\n", "<p data-ntype=\"1\" data-mtype=\"0\"><span class=\"node node--expand\" data-ntype=\"23\" data-mtype=\"2\"><span class=\"marker\">`</span><code data-ntype=\"23\" data-mtype=\"2\">foo</code><span class=\"marker\">`</span></span><span class=\"newline\">\n\n</span></p>"},
	{"5", "**foo**\n", "<p data-ntype=\"1\" data-mtype=\"0\"><span class=\"node node--expand\" data-ntype=\"18\" data-mtype=\"2\"><span class=\"marker\" data-cso=\"2\" data-ceo=\"2\">**</span><strong data-ntype=\"18\" data-mtype=\"2\">foo</strong><span class=\"marker\">**</span></span><span class=\"newline\">\n\n</span></p>"},
	{"4", "_foo_\n", "<p data-ntype=\"1\" data-mtype=\"0\"><span class=\"node node--expand\" data-ntype=\"13\" data-mtype=\"2\"><span class=\"marker\">_</span><em data-ntype=\"13\" data-mtype=\"2\">foo</em><span class=\"marker\">_</span></span><span class=\"newline\">\n\n</span></p>"},
	{"3", "> foo\n", "<div class=\"node node--block node--expand\" data-ntype=\"5\" data-mtype=\"1\"><span class=\"marker\" data-cso=\"2\" data-ceo=\"2\">&gt; </span><blockquote data-ntype=\"5\" data-mtype=\"1\"><p data-ntype=\"1\" data-mtype=\"0\">foo<span class=\"newline\">\n\n</span></p></blockquote></div>"},
	{"2", "## foo\n", "<h2 class=\"node node--expand\" data-ntype=\"2\" data-mtype=\"0\"><span class=\"marker\" data-cso=\"2\" data-ceo=\"2\">## </span>foo</h2>"},
	{"1", "*foo*\n", "<p data-ntype=\"1\" data-mtype=\"0\"><span class=\"node node--expand\" data-ntype=\"13\" data-mtype=\"2\"><span class=\"marker\">*</span><em data-ntype=\"13\" data-mtype=\"2\">foo</em><span class=\"marker\">*</span></span><span class=\"newline\">\n\n</span></p>"},
	{"0", "foo\n", "<p data-ntype=\"1\" data-mtype=\"0\">foo<span class=\"newline\">\n\n</span></p>"},
}

func TestVditorRenderer(t *testing.T) {
	luteEngine := lute.New()

	for _, test := range vditorRendererTests {
		html, err := luteEngine.RenderVditorDOM(test.from, 2, 2)
		if nil != err {
			t.Fatalf("unexpected: %s", err)
		}

		if test.to != html {
			t.Fatalf("test case [%s] failed\nexpected\n\t%q\ngot\n\t%q\noriginal markdown text\n\t%q", test.name, test.to, html, test.from)
		}
	}
}

func TestVditorNewline(t *testing.T) {
	luteEngine := lute.New()

	html, err := luteEngine.VditorNewline(1, nil)
	if nil != err {
		t.Fatalf("unexpected: %s", err)
	}
	expected := "<p data-ntype=\"1\" data-mtype=\"0\"><br><span class=\"newline\">\n</span><span class=\"newline\">\n</span></p>"
	if expected != html {
		t.Fatalf("vditor newline failed\nexpected\n\t%q\ngot\n\t%q", expected, html)
	}
}
