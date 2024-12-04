import { Button, Card, Classes, Checkbox, H5, Navbar, NavbarGroup, NavbarHeading, NavbarDivider, Alignment, Icon, Divider, Tag } from "@blueprintjs/core";
import React from "react";
import { useNavigate } from "react-router-dom";
import { useGetRole } from "../hooks/useAccountInfo.ts";
import { Role } from "./InfoPage.tsx";

interface LandingBarProps {
    showAccount?: boolean;
    showHome?: boolean;
}

export const LandingNavbar = ({ showAccount = true, showHome = true }: LandingBarProps) => {
    const navigate = useNavigate();
    const role = useGetRole();

    const handleClick = () => { showAccount ? navigate('/account') : navigate('/') };

    return (
        <Navbar fixedToTop={true}>
            <Navbar.Group>
                <Icon icon="graph" size={25} />
                <Navbar.Divider />
                <Navbar.Heading><h2>College User Testing</h2></Navbar.Heading>
            </Navbar.Group>
            <Navbar.Group align={Alignment.RIGHT}>
                {showHome && <Navbar.Divider />}
                {showHome && (showAccount ? <Button className="bp5-minimal" icon="user" text="Account" onClick={handleClick} /> : <Button className="bp5-minimal" icon="home" text="Home" onClick={handleClick} />)}
                {(!showAccount && role === Role.Admin) && <Tag round icon="user" large intent="success">Admin</Tag>}
                {(!showAccount && role === Role.Founder) && <Tag round icon="user" large intent="primary" >Founder</Tag>}
                {(!showAccount && role === Role.User) && <Tag round icon="user" large >Tester</Tag>}
            </Navbar.Group>
        </Navbar>
    );
}