// this somehow has to be public
var exports = window;
(function() { 'use strict';

  var jsRunner = function() {
    var PROMPT = "> " ;
    var log_output = PROMPT;
    var ERR_OUTPUT = "Warnings: ";
    var err_output = ERR_OUTPUT;
    var original_log = console.log;
    var original_err = console.error;

    function reset_log_output (){ log_output  = PROMPT;}

    function reset_error_output(){ err_output = ERR_OUTPUT;}

    function redirect_err() {
      err_output = err_output + Array.prototype.slice.apply(arguments).join(' ') + "\n"
    };

    function redirect() {
      log_output = log_output + Array.prototype.slice.apply(arguments).join(' ') + "\n"
    };

    function get_log_output(){
      var old = log_output;
      reset_log_output();
      return old;
    }

    function get_error_output(){
      var old = err_output;
      reset_error_output();
      return old;
    }

    return {
      eval: function (js){
        console.log = redirect;
        try {
          window.eval(js);
          console.log = original_log;
          return get_log_output();
        }
        catch(e){
          var re = (get_log_output() + "\n" + e);
          console.log = original_log;
          return re;
        }
      }
    }
  }();

  var pushedEmpty = (window.location.pathname == "/");
  function playground(opts) {

    var hello = $(opts.inputEl).text();
    $(opts.inputEl).text("");
    var input = CodeMirror(
        $(opts.inputEl)[0],
        {
          mode:'text/x-ocaml',
          lineNumbers: true,
          lineWrapping: true,
          styleActiveLine:true,
          theme: "solarized"
        });
    input.setSize(null,"100%");
    input.setValue(hello);
    var output = CodeMirror(
        $(opts.outputEl)[0],
        {
          mode : 'javascript',
          readOnly: true,
          lineWrapping: true,
          lineNumbers: true,
          theme: "monokai"
        });
    output.setSize(null,"100%");
    if(opts.style['font-size']) {
      $(".CodeMirror").css("font-size", opts.style['font-size']);
    }

    function onInputChanges() {
      if (pushedEmpty) return;
      pushedEmpty = true;
      $(opts.shareURLEl).hide();
      window.history.pushState(null, "", "/");
    }

    function loading() {
      output.setValue('Waiting for remote server...');
    }

    function runner(rType) {
      var runningId = 0;
      return function() {
        loading();
        runningId++;
        var currId = runningId;
        $.ajax("/compile?type="+rType, {
          data: input.getValue(),
          type: "POST",
          dataType: "text",
          beforeSend: loading,
          complete: function(xhr) {
            if (runningId != currId) return;
            if (xhr.status == 200) {
              if (rType === "to_js") {
                output.setValue(jsRunner.eval(xhr.responseText));
              } else {
                output.setValue(xhr.responseText);
              }
              return;
            }
            if (xhr.status == 400) {
              output.setValue(xhr.responseText);
              return;
            }
            alert("Server error; try again.\nError: " + xhr.responseText);
          }
        });
      };
    }

    function fmt() {
      $.ajax("/format", {
        data: input.getValue(),
        type: "POST",
        dataType: "text",
        beforeSend: loading,
        complete: function(xhr) {
          if (xhr.status == 200) {
            input.setValue(xhr.responseText);
            output.setValue("");
            return;
          }
          if (xhr.status == 400) {
            output.setValue(xhr.responseText);
            return;
          }
          alert("Server error; try again.\nError: " + xhr.responseText);
        }
      });
    }

    function origin(href) {
      return (""+href).split("/").slice(0, 3).join("/");
    }
    var shareURL = $(opts.shareURLEl).hide();
    var sharing = false;
    function share() {
      if (sharing) return;
      sharing = true;
      var sharingData = input.getValue();
      $.ajax("/share", {
        processData: false,
        data: sharingData,
        type: "POST",
        complete: function(xhr) {
          sharing = false;
          if (xhr.status != 200) {
            alert("Server error; try again.");
            return;
          }
          var path = "/p/" + xhr.responseText;
          var url = origin(window.location) + path;
          shareURL.show().val(url).focus().select();
          var historyData = {"code": sharingData};
          window.history.pushState(historyData, "", path);
          pushedEmpty = false;
        }
      });
    }

    function embed() {
      if (sharing) return;
      sharing = true;
      var sharingData = input.getValue();
      $.ajax("/share", {
        processData: false,
        data: sharingData,
        type: "POST",
        complete: function(xhr) {
          sharing = false;
          if (xhr.status != 200) {
            alert("Server error; try again.");
            return;
          }
          var path = "/p/" + xhr.responseText;
          var url = origin(window.location) + path;
          var historyData = {"code": sharingData};
          window.history.pushState(historyData, "", path);
          pushedEmpty = false;
          $(opts.embedPopEl).show();
          $(opts.embedOutEl).val(embedHTML(xhr.responseText));
        }
      });
    }

    function embedHTML(id) {
      var attributes = {
        "data-id": id,
        "data-host": window.location.host,
        "data-height": "256",
        "data-font-size": "12px",
        "class": "replay"
      };
      var output = "<p ";
      for (var k in attributes) {
        output += k + '="' + attributes[k] + '" ';
      }
      output += ">\n";
      output += "<!--- Snippet Generated Through Replay\n";
      output += input.getValue() + "\n";
      output += "-->\n";
      output += '</p> <script async src="' + origin(window.location) + '/static/ei.js"></script>';
      return output;
    }

    input.on("changes", onInputChanges);
    $(opts.shareEl).click(share);
    $(opts.runEl).click(runner("to_run"));
    $(opts.fmtEl).click(fmt);
    $(opts.embedEl).click(embed);
    $(window).click(function( event ) {
      if(event.target===$(opts.embedPopEl)[0]) {
        $(opts.embedPopEl).hide();
      }
    });
  }

  window.playground = playground;
})();
