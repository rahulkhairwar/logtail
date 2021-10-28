
function TestTimeout() {
    setTimeout(
        () => {
            console.log("inside setTimeout");
        },
        2000
    );

    console.log("outside setTimeout");
}

export default TestTimeout;