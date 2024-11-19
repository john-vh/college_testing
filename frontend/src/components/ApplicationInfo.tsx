import { Button, Card, CardList, Collapse, FormGroup, H3, H5, Icon, InputGroup, TextArea } from "@blueprintjs/core";
import useAccountInfo from "../hooks/useAccountInfo.ts";
import React, { useEffect, useMemo, useState } from "react";
import { useApplicationInfo } from "../hooks/useApplicationInfo.ts";
import { useBusinessIds } from "../hooks/useBusinessIds.ts";
import { useAcceptApplication } from "../hooks/useAcceptApplication.ts";
import { useRejectApplication } from "../hooks/useRejectApplication.ts";

export const ApplicationInfo = () => {
    const [isOpen, setIsOpen] = useState<boolean[]>([]);
    const data = useApplicationInfo();
    const acceptApplication = useAcceptApplication();
    const rejectApplication = useRejectApplication();

    const toggleCollapse = (index) => {
        setIsOpen((prevState) => ({
            ...prevState,
            [index]: !prevState[index],
        }));
    };

    const handleAcceptApplication = (entry, index) => {
        acceptApplication({ entry, index });
    }

    const handleRejectApplication = (entry, index) => {
        rejectApplication({ entry, index });
    }

    return (
        <div className="content">
            <CardList>
                {data.map((entry, id) => (
                    <Card key={id} style={{ display: "flex", flexDirection: "column" }}>
                        <Button onClick={() => toggleCollapse(id)}>Button 1</Button>
                        <Collapse isOpen={isOpen[id]}>
                            {entry.applications.map((application, index) =>
                                <Card interactive={true} >
                                    <H3>Application Information</H3>
                                    <FormGroup inline label="Posting"
                                        labelFor="posting" >
                                        <InputGroup id={"posting" + index} readOnly />
                                    </FormGroup>
                                    <FormGroup label="Name"
                                        labelFor="name" >
                                        <InputGroup id="name" placeholder={application.user.name} readOnly />
                                    </FormGroup>
                                    <FormGroup label="Email"
                                        labelFor="email" >
                                        <InputGroup id="email" placeholder={application.user.email} readOnly />
                                    </FormGroup>
                                    <FormGroup label="School"
                                        labelFor="school" >
                                        <InputGroup id="school" readOnly />
                                    </FormGroup>
                                    <FormGroup label="Notes"
                                        labelFor="notes" >
                                        <TextArea id="email" defaultValue={application.notes} readOnly fill />
                                    </FormGroup>
                                    <Button onClick={() => handleAcceptApplication(entry, index)}>Accept applicant</Button>
                                    <Button onClick={() => handleRejectApplication(entry, index)} style={{ margin: "10px" }}>Reject applicant</Button>
                                </Card>
                            )}
                        </Collapse>
                    </Card>
                ))}
            </CardList>
        </div>);
}