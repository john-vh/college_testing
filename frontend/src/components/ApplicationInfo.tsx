import { Button, Card, CardList, Collapse, FormGroup, H3, H5, Icon, InputGroup, Label, TextArea } from "@blueprintjs/core";
import useAccountInfo from "../hooks/useAccountInfo.ts";
import React, { useEffect, useMemo, useState } from "react";
import { useApplicationInfo } from "../hooks/useApplicationInfo.ts";
import { useBusinessIds } from "../hooks/useBusinessIds.ts";
import { useAcceptApplication } from "../hooks/useAcceptApplication.ts";
import { useRejectApplication } from "../hooks/useRejectApplication.ts";
import { usePostingNames } from "../hooks/usePostingIds.ts";

export const ApplicationInfo = () => {
    const [isOpen, setIsOpen] = useState<boolean[]>([]);
    const data = useApplicationInfo();
    const names = usePostingNames();
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
        <div className="content2">
            <CardList>
                {data.map((entry, id) => (
                    <div style={{ display: "flex", flexDirection: "column" }}>
                        <Card key={id}>
                            <div style={{ display: "flex", alignItems: "center", justifyContent: "space-between" }}>
                                <H5>{names.get(entry.post_id) || 'Unknown Posting'}</H5>
                                {entry.applications.length > 0 && <Button
                                    onClick={() => toggleCollapse(id)}
                                    rightIcon={isOpen[id] ? "chevron-up" : "chevron-down"}>
                                    {isOpen[id] ? 'Hide Applications' : 'Show Applications'}
                                </Button>}
                            </div>
                            <Collapse isOpen={isOpen[id]}>
                                {entry.applications.map((application, index) =>
                                    <div style={{ paddingTop: "20px" }}>
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
                                    </div>
                                )}
                            </Collapse>
                        </Card>
                    </div>
                ))
                }
            </CardList >
        </div >);
}