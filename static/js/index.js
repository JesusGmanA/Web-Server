$("#create-form").submit(function(event){
    var data = new FormData($("#create-form")[0]);
    event.preventDefault(); //So the page doesn't refresh.
    $.ajax({
        url: "createStudentGrade",
        data: data,
        enctype: 'multipart/form-data',
        contentType: false,
        cache: false,
        processData: false,
        method: 'POST',
        success: function(response){
            activeResponseView(response)
        },
        error: function(response, request){
            console.log("not good")
        }
    })
});

$("#student-avg").submit(function(event){
    var data = new FormData($("#student-avg")[0]);
    console.log("zup")
    event.preventDefault(); //So the page doesn't refresh.
    $.ajax({
        url: "studentAvgInfo",
        data: data,
        enctype: 'multipart/form-data',
        contentType: false,
        cache: false,
        processData: false,
        method: 'POST',
        success: function(response){
            activeResponseView(response)
        },
        error: function(response, request){
            console.log("not good")
        }
    })
});

$("#class-score").submit(function(event){
    var data = new FormData($("#class-score")[0]);
    event.preventDefault(); //So the page doesn't refresh.
    $.ajax({
        url: "classAvgInfo",
        data: data,
        enctype: 'multipart/form-data',
        contentType: false,
        cache: false,
        processData: false,
        method: 'POST',
        success: function(response){
            activeResponseView(response)
        },
        error: function(response, request){
            console.log("not good")
        }
    })
});

function disableMainView(form){
    let formInvalid = false
    $(form+' input').each(function() {
        if ($(this).val() === '') {
          formInvalid = true;
        }
      });
    if (formInvalid){
      alert('One or more fields are empty. Please fill up all fields');
      return false
    }
    else {
        let activeView = document.getElementById('form-request')
    activeView.style.display = (activeView.style.display === 'none'?'block':'none')
        return true;
    }
}

function activeResponseView(response){
    let hiddenView = document.getElementById('message-view');
    hiddenView.style.display = (hiddenView.style.display === 'none'?'block':'none')
    let replyMessage = document.getElementById('reply');
    replyMessage.innerHTML = "" + response;
}