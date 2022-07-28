let GLOBAL_IS_HEADER = false;                       // 鼠标是否位于Header部分
let GLOBAL_IS_BOX = false;                          // 鼠标是否位于课程容器部分
let GLOBAL_NOW_WEEK = -1;                           // 鼠标当前悬停的周次
let GLOBAL_NOW_SECTION = -1;                        // 鼠标当前悬停的节次
let GLOBAL_HINT_TIMEOUT = -1;                       // 全局提示延迟ID


let GLOBAL_STATIC_FUNC = {
    // 禁用HTML页面选中
    BlockHtmlSelect: function () {
        document.body.οncοntextmenu=document.body.οndragstart= document.body.onselectstart=document.body.onbeforecopy=function(){return false;};
        document.body.οnselect=document.body.οncοpy=document.body.οnmοuseup=function(){document.selection.empty();};
    },
    // 初始化显示状态
    ComponentViewInit: function () {
        $("#plans").hide();
        $("#waits").hide();
    }
}

// 页面初始化
$(function () {
    GLOBAL_STATIC_FUNC.BlockHtmlSelect();
    GLOBAL_STATIC_FUNC.ComponentViewInit();
});

// 开启课位可拖动
$(function() {
    $( "ul", ".course-box" ).draggable({
        cancel: "a.ui-icon", // 点击一个图标不会启动拖拽
        revert: "invalid", // 当未被放置时，条目会还原回它的初始位置
        containment: "document",
        helper: "clone",
        cursor: "move",
        zIndex: 9999999,
        appendTo: "#draggable-parent",
    });
});

// 课位拖动详情
$(".course-box").droppable({
    accept: ".course-box > ul",
    activeClass: "ui-state-highlight",
    drop: function( event, ui ) {
        let slfWeek = ui.draggable.attr("week");
        let slfSection = ui.draggable.attr("section");
        let targetWeek = $(this).attr("week");
        let targetSection = $(this).attr("section");
        if ((slfWeek === targetWeek && slfSection === targetSection) === false) {
            if (!$(this).parent().hasClass("box-notallow")) {
                Apis.MoveSection(ui.draggable, $(this), slfWeek, slfSection, targetWeek, targetSection)
            } else {
                hint("课位冲突", "无法将课程调整到该位置。");
            }
        }
    },
    activate: function (event, ui) {

    }
});


// 开启课程可拖动
$(function() {
    $( "li", ".courses" ).draggable({
        cancel: "a.ui-icon", // 点击一个图标不会启动拖拽
        revert: "invalid", // 当未被放置时，条目会还原回它的初始位置
        containment: "document",
        helper: "clone",
        cursor: "move",
        zIndex: 9999999,
        appendTo: "#draggable-parent",
    });
});

// 课程拖动详情
$(".courses").droppable({
    accept: ".courses > li",
    activeClass: "ui-state-highlight",
    drop: function( event, ui ) {
        let parentWeek = ui.draggable.parent().parent().attr("week");
        let parentSection = ui.draggable.parent().parent().attr("section");
        let targetWeek = $(this).parent().attr("week");
        let targetSection = $(this).parent().attr("section");
        if ((parentWeek === targetWeek && parentSection === targetSection) === false) {
            $(".moreinfo-show").removeClass("moreinfo-show")
            $(".same").removeClass("same");
            if (!$(this).parent().hasClass("box-notallow")) {
                Apis.PlanCourseMove(ui.draggable, $(this), parentWeek, parentSection, targetWeek, targetSection);
            } else {
                hint("课位冲突", "无法将课程调整到该位置。");
            }
        }
        $(".course-item").mouseleave(MouseLeaveCourseListener);
        $(".course-item").mouseenter(MouseEnterCourseListener);
        $(".cause-show").removeClass("cause-show");

        $(".box-notallow").removeClass("box-notallow");
        $(".box-notallow-slf").removeClass("box-notallow-slf");
    },
    activate: function (event, ui) {
        $(".course-item").unbind("mouseleave");
        $(".course-item").unbind("mouseenter");
    }
});

// 提示信息
function hint(title, content) {
    $("#hint").addClass("hind-show");
    $("#hint-title").text(title);
    $("#hint-content").text(content);
    clearTimeout(GLOBAL_HINT_TIMEOUT);
    GLOBAL_HINT_TIMEOUT = setTimeout(function () {
        $("#hint").removeClass("hind-show");
    }, 3000)
}

// 隐藏提示信息
$("#hint").click(function () {
    $("#hint").removeClass("hind-show");
    clearTimeout(GLOBAL_HINT_TIMEOUT);
});
