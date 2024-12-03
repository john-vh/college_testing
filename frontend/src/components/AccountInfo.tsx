import { Button, Card, Checkbox, FormGroup, H3, H5, Icon, InputGroup } from "@blueprintjs/core";
import useAccountInfo from "../hooks/useAccountInfo.ts";
import React, { useEffect } from "react";
import { LandingNavbar } from "../components/LandingNavbar.tsx";
import { useIsFounder } from "../hooks/useBusinessInfo.ts";

export const AccountInfo = () => {

    const data = useAccountInfo();
    const isAdmin = data?.roles.includes("admin");
    const isUser = data?.roles.includes("user");
    const isFounder = useIsFounder();

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
                        <Checkbox checked={isAdmin}>Admin</Checkbox>
                        <Checkbox checked={isUser}>User</Checkbox>
                        <Checkbox checked={isFounder}>Founder</Checkbox>
                    </FormGroup>
                    <Button>Save changes</Button>
                </Card>
            </div>
        );
    }
}

export default AccountInfo;