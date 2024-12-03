import './App.css';
import '@blueprintjs/core/lib/css/blueprint.css';
import { Button, Card, Classes, Checkbox, H5, Navbar, NavbarGroup, NavbarHeading, NavbarDivider, Alignment, Icon, Divider } from "@blueprintjs/core";
import React, { useState } from 'react';
import { LandingNavbar } from './components/LandingNavbar.tsx';
import { Landing } from './containers/Landing.tsx';
import { BrowserRouter as Router, Route, Routes } from 'react-router-dom';
import Posting from './containers/Posting.tsx';
import Account from './containers/Account.tsx';



const App = () => {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<Landing />} />
        <Route path="/posting/:id" element={<Posting />} />
        <Route path="/account" element={<Account />} />
      </Routes>
    </Router>
  );
}
export default App;
