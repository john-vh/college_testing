import { Button, Card, Classes, Checkbox, H5, Navbar, NavbarGroup, NavbarHeading, NavbarDivider, Alignment, Icon, Divider, Tag, Colors } from "@blueprintjs/core";
import React from "react";
import { useNavigate } from "react-router-dom";
import { useGetRole } from "../hooks/useAccountInfo.ts";
import { Role } from "./InfoPage.tsx";

interface LandingBarProps {
    showAccount?: boolean;
    showHome?: boolean;
    handleDialog: (boolean) => void;
}

export const LandingNavbar = ({ showAccount = true, showHome = true, handleDialog }: LandingBarProps) => {
    const navigate = useNavigate();
    const role = useGetRole();

    const handleClick = () => { showAccount ? navigate('/account') : navigate('/') };

    return (
        <Navbar style={{ background: Colors.GOLD4 }} fixedToTop={true}>
            <Navbar.Group>
                <Icon icon="generate" size={25} />
                <Navbar.Divider />
                <Navbar.Heading><h2>TestHive: College Startup Testing</h2></Navbar.Heading>
            </Navbar.Group>
            <Navbar.Group align={Alignment.RIGHT}>
                <Button icon={<Icon icon="help" color={Colors.BLACK} />} minimal onClick={() => handleDialog(true)} />
                {showHome && <Navbar.Divider />}
                {showHome && (showAccount ? <Button className="bp5-minimal" icon={<Icon icon="user" color={Colors.BLACK} />} text="Account" onClick={handleClick} />
                    : <Button className="bp5-minimal" icon={<Icon icon="home" color={Colors.BLACK} />} text="Home" onClick={handleClick} />)}
                {(!showAccount && role === Role.Admin) && <Tag round icon="user" large intent="success">Admin</Tag>}
                {(!showAccount && role === Role.Founder) && <Tag round icon="user" large intent="primary" >Founder</Tag>}
                {(!showAccount && role === Role.User) && <Tag round icon="user" large >Tester</Tag>}
            </Navbar.Group>
        </Navbar>
    );
}