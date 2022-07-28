const API = (() => {
    const _exec = Symbol('Exec');
    class API {
        constructor() {
            this[_exec] = function (method, url, data, success) {
                $.ajax({
                    type: method, url,
                    dataType: "json",
                    data: data,
                    success: success,
                    error: function (e) {
                        hint("请求发起异常", e);
                    }
                });
            }
        }

        MoveSection = function(slfEL, targetEL, week, section, targetWeek, targetSection) {
            this[_exec]("POST", "/v1/plan/section/move", {
                week: week,
                section: section,
                targetWeek: targetWeek,
                targetSection: targetSection,
            }, function (json) {
                if (json.Code === 0) {
                    location.reload();
                } else {
                    hint("提示", json.Error)
                }
            });
        }

        // 当课程被成功拖动的时候触发
        // @param courseElement 被拖动的课程元素标签
        // @param targetContainerElement 欲放入的目标容器
        // @param week 拖动前周次
        // @param section 拖动前节次
        // @param targetWeek 拖动后目标周次
        // @param targetSection 拖动后目标节次
        PlanCourseMove = function (courseElement, targetContainerElement, week, section, targetWeek, targetSection) {
            let courseName = courseElement.attr("course");
            if (!GLOBAL_IS_HEADER) {
                if (courseElement.attr("week") == null) {
                    week = -1;
                    section = -1;
                } else if (targetSection == null && targetWeek == null) {
                    targetSection = -1;
                    targetWeek = -1;
                }
            }

            this[_exec]("POST", "/v1/plan/course/move", {
                courseName: courseName,
                week: week,
                section: section,
                targetWeek: targetWeek,
                targetSection: targetSection,
            }, function (json) {
                if (json.Code === 0) {
                    courseElement.fadeOut(function() {
                        let newCourse = targetContainerElement.append(courseElement).children(':last-child');
                        if (newCourse.attr("week") == null) {
                            let content = ""
                            newCourse.children().each(function () {
                                content = content + $(this).text();
                            });
                            newCourse.html(newCourse.attr("course") + `<div class="moreinfo">` + content + `</div>`);
                        }
                        newCourse.fadeIn();

                        if (targetSection === -1 && targetWeek === -1) {
                            newCourse.removeAttr("week");
                            newCourse.removeAttr("section");
                            newCourse.attr("course", newCourse.attr("course"));
                        } else {
                            newCourse.attr("week", targetWeek);
                            newCourse.attr("section", targetSection);
                        }
                    });
                } else {
                    hint("提示", json.Error)
                }
            });
        }

        // 获取特定课程可排课的位置并展示
        // @param courseName 课程名称
        GetWaitCourseAllowSectionAndShow = function (courseName) {
            this[_exec]("GET", "/v1/plan/course/allows", {
                courseName: courseName
            }, function (json) {
                if (json.Code === 0) {
                    // 获取允许排课课位的数量后遍历所有课位容器
                    let length = json.Data[0] == null ? 0: json.Data[0].length;
                    $(`[isbox="true"]`).each(function () {
                        let same = false;
                        for (var i = 0; i < length; i++) {
                            let weekAndSection = json.Data[0][i];
                            if ($(this).attr("week") === weekAndSection[0] + "" && $(this).attr("section") === weekAndSection[1] + "") {
                                same = true;
                                break;
                            }
                        }
                        if (same) {
                            $(this).removeClass("box-notallow");
                            $(this).children(":last").removeClass("cause-show");
                        }else {
                            $(this).addClass("box-notallow");
                            $(this).children(":last").addClass("cause-show");
                            let content = json.Data[1].Cause[parseInt($(this).attr("week"))][parseInt($(this).attr("section"))];
                            if (content === "") {content = "禁排课位"}
                            $(this).children(":last").html(content);
                        }
                    })

                    if (GLOBAL_IS_BOX === false) {
                        leave();
                    }
                } else {
                    hint("提示", json.Error);
                }
            });

        }

        // 获取存在冲突的课位并显示
        // @param courseName 待排查课程
        // @param week 待排查课程所在周次
        // @param section 待排查课程所在节次
        GetUnallowableSectionAndShow = function (courseName, week, section) {
            this[_exec]("GET", "/v1/plan/course/unallowable", {
                courseName: courseName,
                week: week,
                section: section,
            }, function (json) {
                if (week !== GLOBAL_NOW_WEEK || section !== GLOBAL_NOW_SECTION || GLOBAL_IS_BOX === false) { return; }
                if (json.Code === 0) {
                    for (let i = 0; i < json.Data[0].length; i++) {
                        let weekAndSection = json.Data[0][i];
                        let el = $(`[isbox="true"][week="` + weekAndSection[0] + `"][section="` + weekAndSection[1] + `"]`);
                        if (weekAndSection[0] + "" === week && weekAndSection[1] + "" === section){
                            el.addClass("box-notallow-slf")
                        }else {
                            let content = json.Data[1].Cause[weekAndSection[0]][weekAndSection[1]];
                            if (content === "") {content = "禁排课位"}
                            el.addClass("box-notallow")
                            el.children(":last").addClass("cause-show").html(content);
                        }
                    }
                } else {
                    if (json.Error === "请选择方案"){
                        location.reload();
                    }
                }

                setTimeout(function () {
                    if (GLOBAL_IS_BOX === false) {
                        MouseLeaveCourseListener()
                    }
                }, 100);
            })
        }

        // 新建排课方案
        // @param planName 待新建方案名称
        // @param week 待新建方案上课周次数量
        // @param section 待新建方案课节次
        NewPlan = function(planName, week, section) {
            this[_exec]("POST", "/v1/plan/new", {
                planName: planName,
                week: week,
                section: section,
            }, function (json) {
                if (json.Code === 0) {
                    location.reload();
                } else {
                    hint("提示", json.Error)
                }
            });
        }

        // 切换当前打开的排课方案
        // @param planName 目标方案名称
        SwitchPlan = function (planName) {
            this[_exec]("POST", "/v1/plan/switch", {
                planName: planName,
            }, function (json) {
                if (json.Code === 0) {
                    location.reload();
                } else {
                    hint("提示", json.Error)
                }
            });
        }

        // 对当前打开的排课方案进行自动排课
        Auto = function () {
            this[_exec]("POST", "/v1/plan/auto", {}, function (json) {
                if (json.Code === 0) {
                    location.reload();
                } else {
                    hint("提示", json.Error)
                }
            });
        }

        // 对当前打开的排课方案进行冲突优化
        Optimize = function () {
            this[_exec]("POST", "/v1/plan/optimize", {}, function (json) {
                if (json.Code === 0) {
                    location.reload();
                } else {
                    hint("提示", json.Error)
                }
            });
        }

        // 导入排课数据
        // @param formData 导入文件数据
        Import = function (formData) {
            $.ajax({
                type: "POST", url: "/v1/plan/import",
                dataType: "json",
                contentType: false,
                processData: false,
                data : formData,
                success: function (json) {
                    if (json.Code === 0) {
                        location.reload();
                    } else {
                        hint("提示", json.Error)
                    }
                },
                error: function (e) {
                    hint("提示", e);
                }
            })
        }
    }
    return API;
})();
const Apis = new API();
