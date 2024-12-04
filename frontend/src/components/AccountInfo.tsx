import { Button, Card, Checkbox, FormGroup, H3, H5, Icon, InputGroup } from "@blueprintjs/core";
import useAccountInfo, { useGetRole } from "../hooks/useAccountInfo.ts";
import React, { useEffect } from "react";
import { LandingNavbar } from "../components/LandingNavbar.tsx";
import { useIsFounder } from "../hooks/useBusinessInfo.ts";
import { Role } from "./InfoPage.tsx";

export const AccountInfo = () => {

    const data = useAccountInfo();
    const role = useGetRole();

    if (data != null) {
        return (
            <div className="content">
                <Card interactive={true} >
                    <H3>Account Information</H3>
                    <FormGroup label="Name"
                        labelFor="name" >
                        <InputGroup id="name" defaultValue={data.name} />
                    </FormGroup>
                    <FormGroup label="Email"
                        labelFor="email" >
                        <InputGroup id="email" defaultValue={data.email} />
                    </FormGroup>
                    <FormGroup label="Roles"
                        labelFor="roles" >
                        <Checkbox checked={true}>User</Checkbox>
                        <Checkbox checked={role === Role.Founder || role === Role.Admin}>Founder</Checkbox>
                        <Checkbox checked={role === Role.Admin}>Admin</Checkbox>
                    </FormGroup>
                    <Button>Logout</Button>
                </Card>
            </div>
        );
    }
}

export default AccountInfo;