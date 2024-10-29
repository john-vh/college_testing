import logo from './logo.svg';
import './App.css';
import '@blueprintjs/core/lib/css/blueprint.css';
import { Button, Card, Classes, Checkbox, H5, Navbar, NavbarGroup, NavbarHeading, NavbarDivider, Alignment, Icon, Divider } from "@blueprintjs/core";
import React, { useState } from 'react';

const App = () => {
  return (
    <div>
      <Navbar fixedToTop={true}>
        <Navbar.Group>
          <Icon icon="graph" size={25} />
          <Navbar.Divider />
          <Navbar.Heading><h2>College User Testing</h2></Navbar.Heading>
        </Navbar.Group>
        <Navbar.Group align={Alignment.RIGHT}>
          <Navbar.Divider />
          <Button className="bp5-minimal" icon="user" text="Account" />
        </Navbar.Group>
      </Navbar>

      <div className='App'>
        <FilterBar />
        <TestList />
      </div>
    </div>
  );
}

const TestList = () => {
  const items = ["Uber Eats Order Testing", "CSS Formatting Bug Catching", "Free Trial Mobile App Testing", "Test 4", "Test 5", "Test 6", "Test 7", "Test 8"];

  return (
    <div className='Test-list'>
      {items.map((item, index) => (
        <div className='Card'>
          <Card interactive={true} >
            <div className='Flex'>
              <div className='icon-p'>
                <Icon icon="bookmark" size={35}></Icon>
              </div>
              <div>
                <H5>{item}</H5>
                <p>Startup Name</p>
              </div>
            </div>
            <p className="bp5-text-muted">Lorem ipsum dolor sit amet, consectetur adipiscing elit. In iaculis tincidunt est ac sagittis. Nam eleifend lacus in leo vestibulum laoreet convallis ac tortor. Maecenas fermentum in augue vel tristique. Interdum et malesuada fames ac ante ipsum primis in faucibus. Interdum et malesuada fames ac ante ipsum primis in faucibus. Pellentesque gravida sapien vel porta gravida. </p>
            <div className='Flex align-right'>
              <div className="gap">Compensation: $5</div>
              <Button intent="primary">Details</Button>
            </div>
          </Card>
        </div>

      ))}
    </div>
  );
}

const FilterBar = () => {
  const [isOpen, setOpen] = useState(true);

  if (isOpen) {
    return (
      <div className='Filter-bar'>
        <div className="Filter-header">
          <div><strong>Filter Bar</strong></div>
          <Button intent="primary" onClick={() => setOpen(!isOpen)} icon='filter'></Button>
        </div>
        <Checkbox label='Filter 1' />
        <Checkbox label='Filter 2' />
        <Checkbox label='Filter 3' />
      </div>
    );
  }

  else {
    return (
      <div className='hover'>
        <Button intent="primary" onClick={() => setOpen(!isOpen)} icon='filter'></Button>
      </div>

    )
  }
}
export default App;
