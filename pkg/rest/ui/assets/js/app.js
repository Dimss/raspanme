$(document).ready(function () {

    console.log("doc is ready")
    initCategories()

    $("#exam-selector").on("change", (e) => {
        loadQuestions(e.target.value)
    })


});


function initCategories() {
    $.get("/v1/category", function (data, status) {
        $("#exam-selector").empty()
        data.forEach((el) => {
            $("#exam-selector").append(`<option value=${el.ID}>${el.Description}</option>`)
        })
        loadQuestions(data[0].ID)
    });
}

function loadQuestions(catId) {
    $.get(`/v1/question/${catId}`, function (data, status) {
        $("#questions").empty()
        data.forEach((el) => {
            $("#questions").append(`
                 <div class="row">
                    <div class="col" dir="rtl">
                        ${parseQuestionSection(el[0])}
                    </div>
                    <div class="col" dir="ltr">
                        ${parseQuestionSection(el[1])}
                    </div>
                 </div>
                <div class="row">
                    <hr class="mt-1 mb-1"/>
                </div>
                `)
        })
        $("#questions").append(loadRightAnswers(data))
    });
}

function parseQuestionSection(question) {
    let answers = ''
    question.Answer.sort(function (a, b) {
        return (a.AnswerIndex - b.AnswerIndex);
    });
    question.Answer.forEach((answer) => {
        answers += `
                <div class="row">
                    <div class="col-1">
                        <input
                        onclick="(()=>{
                                if ($('#${answer.ID}').is(':checked')){
                                    if (${answer.Right}){
                                        $('#${answer.ID}-text').css('color', 'green');
                                    }else{
                                        $('#${answer.ID}-text').css('color', 'red');    
                                    }
                                }else{
                                    $('#${answer.ID}-text').css('color', '');
                                }
                            })()" 
                        type="checkbox" id="${answer.ID}" data-is-right="${answer.Right}">
                    </div>
                    <div class="col">
                        <label for="${answer.ID}" id="${answer.ID}-text">
                            <span class="badge rounded-pill bg-primary">
                                ${answer.AnswerIndex}
                            </span> 
                            ${answer.Answer}
                        </label>
                    </div>
                </div>
        `
    })

    return `
        <form>
            <div class="form-check">
                <h6><span class="badge bg-secondary">${question.ID}</span> ${question.Question.trim()}</h6>
                ${answers}
            </div>
        </form>
    `
}

function loadRightAnswers(data) {
    let rightAnswers = ''
    data.forEach((el) => {
        el.forEach((question) => {
            question.Answer.forEach((answer) => {
                if (answer.Right) {
                    rightAnswers += `<span class="badge bg-secondary">${question.ID}:${answer.AnswerIndex}</span>`
                }
            })
        })
    })
    return rightAnswers
}
