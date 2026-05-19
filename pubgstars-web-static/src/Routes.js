import React from "react";
import { Route, Switch } from "react-router-dom";
import Home from "./containers/Home";
import Login from "./containers/Login";
import Signup from "./containers/Signup";
import Rules from "./containers/Rules";
import LeaderBoard from "./containers/LeaderBoard";
import Sss from "./containers/Sss";
import About from "./containers/About";
import Contact from "./containers/Contact";
import ResetPassword from "./containers/ResetPassword";
import Balance from "./containers/Balance";
import TransactionLog from "./containers/TransactionLog";
import MyGames from "./containers/MyGames";
import Profile from "./containers/Profile";
import NotFound from "./containers/NotFound";
import AppliedRoute from "./components/AppliedRoute";
import AuthenticatedRoute from "./components/AuthenticatedRoute";
import UnauthenticatedRoute from "./components/UnauthenticatedRoute";


export default ({ childProps }) =>
  <Switch>
    <AppliedRoute path="/" exact component={Home} props={childProps} />
    <UnauthenticatedRoute path="/login" exact component={Login} props={childProps} />
    <UnauthenticatedRoute path="/signup" exact component={Signup} props={childProps} />
    <UnauthenticatedRoute path="/rules" exact component={Rules} props={childProps} />
    <UnauthenticatedRoute path="/sss" exact component={Sss} props={childProps} />
    <UnauthenticatedRoute path="/about" exact component={About} props={childProps} />
    <UnauthenticatedRoute path="/contact" exact component={Contact} props={childProps} />
    <UnauthenticatedRoute path="/login/reset" exact component={ResetPassword} props={childProps}/>
    <AuthenticatedRoute path="/balance" exact component={Balance} props={childProps} />
    <AuthenticatedRoute path="/transactionLog" exact component={TransactionLog} props={childProps} />
    <AuthenticatedRoute path="/leaderboard" exact component={LeaderBoard} props={childProps} />
    <AuthenticatedRoute path="/oyunlarim" exact component={MyGames} props={childProps} />
    <AuthenticatedRoute path="/profile" exact component={Profile} props={childProps} />
    <AuthenticatedRoute path="/app/rules" exact component={Rules} props={childProps} />
    <AuthenticatedRoute path="/app/contact" exact component={Contact} props={childProps} />

    { /* Finally, catch all unmatched routes */ }
    <Route component={NotFound} />
  </Switch>;
