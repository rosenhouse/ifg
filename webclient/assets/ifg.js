var gridster;
var add_new;

$(function(){

var custom_serializer = function($w, wgd) { 
    return { 
        Layout: {
            Col: wgd.col,
            Row: wgd.row,
            SizeX: wgd.size_x,
            SizeY: wgd.size_y
        },
        Content: {
            URL: $w.find("iframe").attr("src")
        }
    };
}
gridster = $(".gridster ul").gridster({
  widget_base_dimensions: [20, 20],
  widget_margins: [2, 2],
  helper: 'clone',
  resize: { enabled: true },
  serialize_params: custom_serializer
}).data('gridster');

var editable = true

var set_editable = function(e) {
    editable = e;
    if (e) {
        $('[role="content"]').hide();
        $('.gs-resize-handle').show();
        gridster.enable();
        $('[role="btn-new-item"]').show();
        $('[role="link-clone"]').show();
        $('[role="controls"]').show();
        $('[role="btn-edit"]').html('Save');
    }
    else {
        $('[role="content"]').show();
        $('.gs-resize-handle').hide();
        gridster.disable();
        $('[role="btn-new-item"]').hide();
        $('[role="link-clone"]').hide();
        $('[role="controls"]').hide();
        $('[role="btn-edit"]').html('Edit');
    }
}

var delete_click = function(widget) {
    return function(event) { 
         gridster.remove_widget(widget);
    };
}

var afterDelay = function(changeAction) {
    var timerid;    
    return function(e) {
        var target = e.target
        var value = $(target).val();
        if($(target).data("lastval")!= value){
            $(target).data("lastval",value);

            clearTimeout(timerid);
            timerid = setTimeout(function() { changeAction(value); }, 500);
        };
    };
}; 

add_new = function(itemSpecs) {
    var iframe = $('<iframe />', {
        role: 'content',
        style: "position: relative; height: 100%; width: 100%; border: none;" });
    iframe.hide();

    var controls = $('<div />', { role: 'controls' });
    var edit_url = $('<input />', { type: 'text', role: 'txt-url', style: 'margin: 1em; width: 80%;' });
    var set_iframe_url = function(url) { iframe.attr("src", url); edit_url.val(url); }
    edit_url.on('input', afterDelay(set_iframe_url));
    
    var button = $('<button />', { text: 'X', style: 'margin: 0.9em; float: right; color:red; width: 2em; height: 2em;' });
    controls.append(edit_url, button);

    var newItem = $('<li />');
    $(button).click(delete_click(newItem));
    newItem.append(iframe, controls);

    if (itemSpecs && itemSpecs.Layout && itemSpecs.Layout.SizeX) {
        var layout = itemSpecs.Layout;
        gridster.add_widget(newItem, layout.SizeX, layout.SizeY, layout.Col, layout.Row);
        set_iframe_url(itemSpecs.Content.URL);
    } else {
        gridster.add_widget(newItem, 25, 15);
        set_iframe_url("http://example.com");
    }

    return set_iframe_url;
};


var upload_state = function() {
    var data = gridster.serialize()

    $.ajax({
        type: 'PUT',
        dataType: 'json',
        url: window.location.pathname + "/data",
        data: JSON.stringify(data)
    }).done(function( data ) {
        console.log(data);
    });
}

var download_state = function() {
    $.ajax({
        type: 'GET',
        dataType: 'json',
        url: window.location.pathname + "/data"
    }).done(function( data ) {
    gridster.remove_all_widgets();
    $.each(data, function( index, value ) {
      add_new(value);
    });
    if ($(data).length) {
        console.log("found data, will not be editable");
        set_editable(false);
    } else {
        console.log("no data, will be editable");
        set_editable(true);
    }
  });
}

$('[role="btn-new-item"]').click(add_new);
$('[role="link-clone"]').attr('href', window.location.href + "/clone");

$('[role="btn-edit"]').on('click', function() {
    set_editable(!editable);
    upload_state();
});


download_state();
});

