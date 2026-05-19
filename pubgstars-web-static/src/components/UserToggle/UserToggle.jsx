import React from "react";

import "./UserToggle.css";

class UserToggle extends React.Component {
    constructor(props, context) {
        super(props, context);
        this.handleClick = this.handleClick.bind(this);
    }

    handleClick(e) {
        e.preventDefault();
        this.props.onClick(e);
    }

    render() {
        return (
            <div className="login-icon" onClick={this.handleClick} ref={this.props.innerRef}>
                <i className="fa fa-user-circle" aria-hidden="true"/>
            </div>
        );
    }
}

export default React.forwardRef((props, ref) => <UserToggle innerRef={ref} {...props}/>);
