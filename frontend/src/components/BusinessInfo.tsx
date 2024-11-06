import { Button, Card, FormGroup, H3, H5, Icon, InputGroup } from "@blueprintjs/core";
import useAccountInfo from "../hooks/useAccountInfo.ts";
import React, { useEffect } from "react";
import { LandingNavbar } from "../components/LandingNavbar.tsx";
import { useBusinessInfo } from "../hooks/useBusinessInfo.ts";

export const BusinessInfo = () => {

    const data = useBusinessInfo()[0];

    if (data != null) {
        return (
            <div className="content">
                <Card interactive={true} >
                    <H3>Account Information</H3>
                    <FormGroup label="Name"
                        labelFor="name" >
                        <InputGroup id="name" placeholder={data.name} />
                    </FormGroup>
                    <FormGroup label="Email"
                        labelFor="email" >
                        <InputGroup id="email" placeholder={data.desc} />
                    </FormGroup>
                    <Button>Save changes</Button>
                </Card>
            </div>
        );
    }
}