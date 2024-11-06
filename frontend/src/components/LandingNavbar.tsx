import { Button, Card, Classes, Checkbox, H5, Navbar, NavbarGroup, NavbarHeading, NavbarDivider, Alignment, Icon, Divider } from "@blueprintjs/core";
import React from "react";
import { useNavigate } from "react-router-dom";

interface LandingBarProps {
    showAccount?: boolean;
}

export const LandingNavbar = ({ showAccount = true }: LandingBarProps) => {
    const navigate = useNavigate();
    const handleClick = () => { showAccount ? navigate(`/account`) : navigate('/'); };

    return (
        <Navbar fixedToTop={true}>
            <Navbar.Group>
                <Icon icon="graph" size={25} />
                <Navbar.Divider />
                <Navbar.Heading><h2>College User Testing</h2></Navbar.Heading>
            </Navbar.Group>
            <Navbar.Group align={Alignment.RIGHT}>
                <Navbar.Divider />
                {showAccount ? <Button className="bp5-minimal" icon="user" text="Account" onClick={handleClick} /> : <Button className="bp5-minimal" icon="home" text="Home" onClick={handleClick} />}

            </Navbar.Group>
        </Navbar>
    );
}