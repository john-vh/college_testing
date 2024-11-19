import { Button, Card, FormGroup, H3, H5, Icon, InputGroup, TextArea } from "@blueprintjs/core";
import useAccountInfo from "../hooks/useAccountInfo.ts";
import React, { useEffect, useState } from "react";
import { LandingNavbar } from "../components/LandingNavbar.tsx";
import { useBusinessInfo } from "../hooks/useBusinessInfo.ts";
import { AddBusiness } from "./AddBusiness.tsx";

export const BusinessInfo = () => {

    const data = useBusinessInfo();
    const [businessAdd, setBusinessAdd] = useState(false);

    const addNewBusiness = () => {
        setBusinessAdd(true); // Change state to show the Add Business page
    };

    if (businessAdd) {
        return <AddBusiness setBusinessAdd={setBusinessAdd} />
    }

    return (
        <div style={{ width: '100%', padding: '20px' }}>
            <Button
                intent="primary"
                large={true}
                style={{ width: "100%", marginBottom: "20px" }}
                onClick={() => addNewBusiness()}
            >
                Add Business
            </Button>
            <div style={{ display: 'flex', flexDirection: 'column', gap: '20px' }}>
                {data != null && data.map((entry) => (
                    <Card interactive={true} >
                        <H3>Business Information</H3>
                        <FormGroup label="Name"
                            labelFor="name" >
                            <InputGroup id="name" defaultValue={entry.name} readOnly />
                        </FormGroup>
                        <FormGroup label="Website"
                            labelFor="website" >
                            <InputGroup id="website" defaultValue={entry.website} readOnly />
                        </FormGroup>
                        <FormGroup label="Description"
                            labelFor="desc" >
                            <TextArea id="desc" placeholder={entry.desc} readOnly fill />
                        </FormGroup>
                        <Button>Manage business</Button>
                    </Card>
                ))}
            </div>

        </div>
    );
}
