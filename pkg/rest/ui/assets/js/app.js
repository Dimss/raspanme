$(document).ready(function () {

    console.log("doc is ready")
    initCategories()

    $("#exam-selector").on("change", (e) => {
        loadQuestions(e.target.value)
    })


});


function initCategories() {
    $.get("/v1/category", function (data, status) {
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
            $("#questions").append(
                `
                 <div class="row">
                    <div class="col" dir="rtl" >
                        ${parseQuestionSection(el[0])}
                    </div>
                    <div class="col" dir="ltr">
                        ${parseQuestionSection(el[1])}
                    </div>
                 </div>
                <div class="row">
                    <hr class="mt-1 mb-1"/>
                </div>
                `
            )

        })
        console.log(data)
    });
}

function parseQuestionSection(question) {
    return `
        <form>
            <div class="form-check">
                <h6>${question.Question.trim()}</h6>
                <div class="row">
                    <div class="col-1">
                        <input type="checkbox" id="${question.Answer[0].ID}">
                    </div>
                    <div class="col">
                        <label for="${question.Answer[0].ID}">${question.Answer[0].Answer}</label>
                    </div>
                </div>
                <div class="row">
                    <div class="col-1">
                        <input type="checkbox" id="${question.Answer[1].ID}">
                    </div>
                    <div class="col">
                        <label for="${question.Answer[1].ID}">${question.Answer[1].Answer}</label>
                    </div>
                </div>
                <div class="row">
                    <div class="col-1">
                        <input type="checkbox" id="${question.Answer[2].ID}">
                    </div>
                    <div class="col">
                        <label for="${question.Answer[2].ID}">${question.Answer[2].Answer}</label>
                    </div>
                </div>
                <div class="row">
                    <div class="col-1">
                        <input type="checkbox" id="${question.Answer[3].ID}">
                    </div>
                    <div class="col">
                        <label for="${question.Answer[3].ID}">${question.Answer[3].Answer}</label>
                    </div>
                </div>
            </div>
        </form>
    `
}