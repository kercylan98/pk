// 鼠标进入课程时触发
let MouseEnterCourseListener = function () {
    $(`[course='` + $(this).attr("course") + `']`).addClass("same");
    $(this).children().each(function () {
        $(this).addClass("moreinfo-show")
    });

    if ($(this).attr("week") == null) {
        Apis.GetWaitCourseAllowSectionAndShow($(this).attr("course"));
        return;
    }
    let week = $(this).attr("week")
    let section = $(this).attr("section")
    GLOBAL_NOW_WEEK = week;
    GLOBAL_NOW_SECTION = section;

    Apis.GetUnallowableSectionAndShow($(this).attr("course"), week, section);
}

// 鼠标离开课程时触发
let MouseLeaveCourseListener = function () {
    $(".same").removeClass("same");
    $(".cause-show").removeClass("cause-show");
    $(this).children().each(function () {
        $(this).removeClass("moreinfo-show")
    });

    $(".box-notallow").removeClass("box-notallow");
    $(".box-notallow-slf").removeClass("box-notallow-slf");
}

// 设置课程容器和待排课区域鼠标进入的时候表示在容器中
$(".course-box, #waits").mouseenter(function () {
    GLOBAL_IS_BOX = true;
});

// 设置课程容器和待排课区域鼠标离开的时候表示不在容器中
$(".course-box, #waits").mouseleave(function () {
    GLOBAL_IS_BOX = false;
})

// 设置鼠标移动进入课程时处理函数
$(".course-item").mouseenter(MouseEnterCourseListener);

// 设置鼠标离开课程时处理函数
$(".course-item").mouseleave(MouseLeaveCourseListener);

// 设置鼠标进入顶部菜单时触发
$(".root-header").mouseenter(function () {
    GLOBAL_IS_HEADER = true;
});

// 设置鼠标离开顶部菜单时触发
$(".root-header").mouseleave(function () {
    GLOBAL_IS_HEADER = false;
});

// 设置body在鼠标抬起的时候处理(解决拖到课表外区域课表视觉效果卡死的问题)
$("body").mouseup(function (e) {
    let inRang = false;
    let courses = $(".courses");
    for (let i = 0; i < courses.length; i++) {
        if(courses[i] === e.target || $.contains(courses[i], e.target)){
            inRang = true;
            break;
        }else {
            inRang = false;
        }
    }
    if (!inRang) {
        MouseLeaveCourseListener();
        $(".course-item").mouseenter(MouseEnterCourseListener);
        $(".course-item").mouseleave(MouseLeaveCourseListener);
    }

});

// 切换当前打开的排课方案按下时触发
$(".plans-item").click(function () {
    Apis.SwitchPlan($(this).text());
});

// 排课方案列表按下触发
$("#plan-name").click(function () {
    $("#plans").toggle();
    $("#waits").hide();
});

// 排课方案列表鼠标离开隐藏
$("#plans").mouseleave(function () {
    $("#plans").hide();
});

// 待排课按下显示
$("#waits-name").click(function () {
    $("#waits").toggle();
    $("#plans").hide();
});

// 新建方案
$("#new-plan").click(function () {
    let name = prompt("请输入方案名称", "");
    if(name === "" || name.length === 0) {
        hint("提示", "方案名称不能为空")
        return
    }

    let week = prompt("请输入上课周数", "5");
    if(week === "" || week.length === 0) {
        hint("提示", "上课周数不能为空")
        return
    }

    let section = prompt("请输入每日课节数", "11");
    if(section === "" || section.length === 0) {
        hint("提示", "每日课节数不能为空")
        return
    }

    Apis.NewPlan(name, week, section);
});

// 优化课表
$("#auto-optimize").click(function () {
    if(confirm('优化课表将会对当前课表进行调整，是否确定？')){
        Apis.Optimize();
    }
    location.reload();
})

// 自动排课
$("#auto-build").click(function () {
    if(confirm('自动排课将会覆盖现有课表，是否确定？')){
        Apis.Auto();
    }
    location.reload();
})

// 上传数据
$("#upload-data").click(function () {
    $("#upload-data-file").click();
});

// 选择上传数据的文件时触发
$("#upload-data-file").change(function(){
    if($(this).val()){
        var formData = new FormData();
        formData.append('file', $('#upload-data-file')[0].files[0]);
        formData.append('planName', $('#plan-name').text());

        Apis.Import(formData);
    }
});

// 年级筛选
$(".select-stage").click(function () {
    let stage = $(this).text();
    if (stage === "全部年级") {
        $("*[stage]").show();
    }else {
        $("*[stage]").hide();
        $("*[stage='" + stage + "']").show();
    }
});

// 导入模板下载
$("#download-template").click(function () {

});