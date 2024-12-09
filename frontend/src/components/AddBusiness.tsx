import {
    Button,
    Card,
    FormGroup,
    H3,
    InputGroup,
    Intent,
    Classes,
    Callout,
    Icon,
    Spinner
} from "@blueprintjs/core";
import React, { useState } from "react";
import { useCreateBusiness, NewBusinessInfo } from "../hooks/useCreateBusiness.ts";

interface AddBusinessProps {
    setBusinessAdd: (value: boolean) => void;
    fetchData: () => Promise<void>;
}

export const AddBusiness = ({ setBusinessAdd, fetchData }: AddBusinessProps) => {
    const { createBusiness } = useCreateBusiness();
    const [isLoading, setIsLoading] = useState(false);
    const [errors, setErrors] = useState<Partial<Record<keyof NewBusinessInfo, string>>>({});

    const delay = (ms: number) => new Promise(resolve => setTimeout(resolve, ms));

    const [userInfo, setUserInfo] = useState<NewBusinessInfo>({
        name: "",
        desc: "",
        website: ""
    });

    const validateForm = (): boolean => {
        const newErrors: Partial<Record<keyof NewBusinessInfo, string>> = {};

        if (!userInfo.name.trim()) {
            newErrors.name = "Business name is required";
        }

        if (!userInfo.desc.trim()) {
            newErrors.desc = "Description is required";
        }

        if (userInfo.website && !userInfo.website.match(/^(https?:\/\/)?([\da-z.-]+)\.([a-z.]{2,6})([/\w .-]*)*\/?$/)) {
            newErrors.website = "Please enter a valid website URL";
        }

        setErrors(newErrors);
        return Object.keys(newErrors).length === 0;
    };

    const handleChangeInfo = (id: keyof NewBusinessInfo, value: string) => {
        setUserInfo((prevInfo) => ({
            ...prevInfo,
            [id]: value
        }));

        // Clear error when user starts typing
        if (errors[id]) {
            setErrors((prev) => ({
                ...prev,
                [id]: undefined
            }));
        }
    };

    const handleSubmit = async () => {
        if (!validateForm()) {
            return;
        }

        try {
            setIsLoading(true);
            await createBusiness(userInfo);
            await delay(100);
            await fetchData();
            setBusinessAdd(false);
        } catch (error) {
            console.error("Failed to create business:", error);
        } finally {
            setIsLoading(false);
        }
    };

    return (
        <div className="content">
            <Card
                style={{
                    width: "100%",
                    maxWidth: "500px",
                    margin: "auto",
                    padding: "2rem"
                }}
            >
                <H3>Add New Business</H3>

                <div style={{ marginTop: "1.5rem" }}>
                    {Object.values(errors).some(error => error) && (
                        <Callout
                            intent={Intent.DANGER}
                            icon="error"
                            style={{ marginBottom: "1rem" }}
                        >
                            Please fix the errors below to continue
                        </Callout>
                    )}

                    <FormGroup
                        label="Business Name"
                        labelFor="name"
                        intent={errors.name ? Intent.DANGER : Intent.NONE}
                        helperText={errors.name}
                        labelInfo="(required)"
                    >
                        <InputGroup
                            id="name"
                            value={userInfo.name}
                            onChange={(e) => handleChangeInfo("name", e.target.value)}
                            placeholder="Business name"
                            intent={errors.name ? Intent.DANGER : Intent.NONE}
                            leftIcon="office"
                        />
                    </FormGroup>

                    <FormGroup
                        label="Description"
                        labelFor="desc"
                        intent={errors.desc ? Intent.DANGER : Intent.NONE}
                        helperText={errors.desc}
                        labelInfo="(required)"
                    >
                        <InputGroup
                            id="desc"
                            value={userInfo.desc}
                            onChange={(e) => handleChangeInfo("desc", e.target.value)}
                            placeholder="Description"
                            intent={errors.desc ? Intent.DANGER : Intent.NONE}
                            leftIcon="annotation"
                        />
                    </FormGroup>

                    <FormGroup
                        label="Website"
                        labelFor="website"
                        intent={errors.website ? Intent.DANGER : Intent.NONE}
                        helperText={errors.website}
                    >
                        <InputGroup
                            id="website"
                            value={userInfo.website}
                            onChange={(e) => handleChangeInfo("website", e.target.value)}
                            placeholder="https://example.com"
                            intent={errors.website ? Intent.DANGER : Intent.NONE}
                            leftIcon="globe-network"
                        />
                    </FormGroup>

                    <div style={{
                        display: "flex",
                        gap: "0.75rem",
                        marginTop: "2rem"
                    }}>
                        <Button
                            onClick={handleSubmit}
                            intent={Intent.PRIMARY}
                            disabled={isLoading}
                            icon={isLoading ? <Spinner size={16} /> : "tick"}
                        >
                            {isLoading ? "Creating..." : "Create Business"}
                        </Button>
                        <Button
                            onClick={() => setBusinessAdd(false)}
                        >
                            Cancel
                        </Button>
                    </div>
                </div>
            </Card>
        </div>
    );
};