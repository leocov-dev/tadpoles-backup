// ==Bookmarklet==
// @name reset-tadpoles-password
// @author Leo Covarrubias
// @script !loadOnce https://cdnjs.cloudflare.com/ajax/libs/jquery/3.5.0/jquery.min.js
// ==/Bookmarklet==
"use strict";

function findAndClick(selector) {
    const element = $(selector);
    if (element.length <= 0) {
        console.error(`Could not find element with selector "${selector}".`);
        return false;
    }
    element[0].click();  // on first index makes click on a tags work
    return true;
}

function sendEnterEvent(element) {
    const enterEvent = jQuery.Event("keypress");
    enterEvent.which = 13; //choose the one you want
    enterEvent.keyCode = 13;
    element.trigger(enterEvent);
}

function dispatchEnterEvent(element) {
    const ke = new KeyboardEvent("keydown", {
        bubbles: true, cancelable: true, keyCode: 13
    });
    element.trigger(ke);
}

function bringToResetForm() {
    console.info('Starting reset-tadpoles-password script');

    if (window.location.hostname !== 'www.tadpoles.com' &&
        window.location.pathname !== '/home_or_work') {
        alert('Will redirect to https://tadpoles.com/home_or_work\n\nAfterward click bookmarklet again.');
        window.location.replace('https://www.tadpoles.com/home_or_work');
    }

    // if the page already contains the expected submit button, just show it and exit
    const submitButton = $('button[type=submit]');
    const requestHeading = $('div h1.tp-heading-text:contains("Request a password reset")');
    const infoBox = $('div.alert.alert-block.alert-info.tp-centered-contents[data-bind="visible: emailRequest.isGmail()"]');

    if (submitButton.length > 0 && requestHeading.length > 0) {
        console.info('Found reset request submit button, showing...');
        submitButton.show();
        infoBox.hide();
        return;
    }

    if (!findAndClick('div.tp-block-half.tp-centered-contents.tp-pointable[data-bind="click: chooseParent"]')) return;

    if (!findAndClick('img.other-login-button.tp-centered[data-bind="click:chooseTadpoles"]')) return;

    if (!findAndClick('form div.control-group div.controls a[data-bind="click: chooseForgot"]')) return;

    const emailInput = $('form.tp-left-contents input[type=text].pull-left');
    if (emailInput.length <= 0) {
        console.error('Failed to find email input');
        return;
    }

    emailInput.onkeypress((event) => {
        console.info('onkeypress!');
        console.info(event);
        submitButton.show();
    });

    console.log('Completed reset-tadpoles-password script successfully!');


}

// execute bookmarklet
bringToResetForm();
